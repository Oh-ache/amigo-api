package mqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// Consumer is the interface for consuming tasks
type Consumer interface {
	RegisterHandler(name string, handler Handler)
	Start(ctx context.Context) error
	Stop() error
}

// RedisConsumer implements Consumer using asynq
type RedisConsumer struct {
	client   *asynq.Client
	server   *asynq.Server
	handlers map[string]Handler
	config   *QueueConfig
	stopChan chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

// NewRedisConsumer creates a new RedisConsumer
func NewRedisConsumer(redisOpt *redis.Options, config *QueueConfig) *RedisConsumer {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisOpt.Addr, Password: redisOpt.Password, DB: redisOpt.DB})

	// Create asynq server config
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisOpt.Addr, Password: redisOpt.Password, DB: redisOpt.DB},
		asynq.Config{
			Queues: map[string]int{
				"default":   6,
				"critical":  10,
				"low":       1,
			},
			Concurrency: config.Concurrency,
		},
	)

	return &RedisConsumer{
		client:   client,
		server:   srv,
		handlers: make(map[string]Handler),
		config:   config,
		stopChan: make(chan struct{}),
	}
}

// RegisterHandler registers a task handler
func (c *RedisConsumer) RegisterHandler(name string, handler Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[name] = handler
}

// Start starts the consumer
func (c *RedisConsumer) Start(ctx context.Context) error {
	mux := c.createMux()

	logx.Infof("Starting mqueue consumer with %d workers", c.config.Concurrency)

	go func() {
		if err := c.server.Run(mux); err != nil {
			logx.Errorf("Mqueue server error: %v", err)
		}
	}()

	// Keep running until stopped
	go func() {
		<-c.stopChan
		c.server.Stop()
	}()

	return nil
}

// Stop stops the consumer
func (c *RedisConsumer) Stop() error {
	close(c.stopChan)
	c.wg.Wait()
	c.client.Close()
	logx.Info("Mqueue consumer stopped")
	return nil
}

// createMux creates the asynq serve mux with registered handlers
func (c *RedisConsumer) createMux() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	c.mu.RLock()
	defer c.mu.RUnlock()

	for name, handler := range c.handlers {
		handlerName := name
		h := handler
		mux.HandleFunc(handlerName, func(ctx context.Context, t *asynq.Task) error {
			return c.handleTask(ctx, t, h)
		})
	}

	return mux
}

// handleTask processes a single task
func (c *RedisConsumer) handleTask(ctx context.Context, t *asynq.Task, handler Handler) error {
	startTime := time.Now()

	// Parse task payload
	var task Task
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		logx.Errorf("Unmarshal task error: %v", err)
		return fmt.Errorf("unmarshal task error: %w", err)
	}

	// Set timeout if specified
	if task.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(task.Timeout)*time.Second)
		defer cancel()
	}

	// Execute handler
	err := handler.Handle(ctx, &task)

	duration := time.Since(startTime).Milliseconds()

	// Record result
	c.recordResult(&task, err, duration)

	if err != nil {
		logx.Errorf("Handle task %s error: %v, duration: %dms", task.ID, err, duration)
		return err
	}

	logx.Infof("Handle task %s success, duration: %dms", task.ID, duration)
	return nil
}

// recordResult saves task execution result
func (c *RedisConsumer) recordResult(task *Task, err error, duration int64) {
	status := StatusCompleted
	var errMsg string
	if err != nil {
		status = StatusFailed
		errMsg = err.Error()
	}

	result := &TaskResult{
		TaskID:     task.ID,
		Status:     status,
		Error:      errMsg,
		Duration:   duration,
		FinishedAt: time.Now().Unix(),
	}

	resultData, _ := json.Marshal(result)
	resultKey := fmt.Sprintf("mqueue:result:%s", task.ID)
	rdb := redis.NewClient(c.config.RedisOpt)
	defer rdb.Close()

	// Store result with 7 day expiration
	rdb.SetEx(context.Background(), resultKey, string(resultData), 7*24*time.Hour)
}

// GetTaskResult retrieves task execution result
func (c *RedisConsumer) GetTaskResult(ctx context.Context, taskID string) (*TaskResult, error) {
	resultKey := fmt.Sprintf("mqueue:result:%s", taskID)
	rdb := redis.NewClient(c.config.RedisOpt)
	defer rdb.Close()

	resultData, err := rdb.Get(ctx, resultKey).Bytes()
	if err != nil {
		return nil, fmt.Errorf("get task result error: %w", err)
	}

	var result TaskResult
	if err := json.Unmarshal(resultData, &result); err != nil {
		return nil, fmt.Errorf("unmarshal task result error: %w", err)
	}

	return &result, nil
}

