package mqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// Producer is the interface for enqueuing tasks
type Producer interface {
	Enqueue(ctx context.Context, task *Task) (string, error)
	EnqueueWithPriority(ctx context.Context, task *Task, priority Priority) (string, error)
	EnqueueDelayed(ctx context.Context, task *Task, delay time.Duration) (string, error)
	EnqueueToQueue(ctx context.Context, task *Task, queueName string) (string, error)
	BatchEnqueue(ctx context.Context, tasks []*Task) ([]string, error)
	CancelTask(ctx context.Context, taskID string) error
}

// RedisProducer implements Producer using asynq
type RedisProducer struct {
	client *asynq.Client
	config *QueueConfig
}

// NewRedisProducer creates a new RedisProducer
func NewRedisProducer(redisOpt *redis.Options, config *QueueConfig) *RedisProducer {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisOpt.Addr, Password: redisOpt.Password, DB: redisOpt.DB})
	return &RedisProducer{
		client: client,
		config: config,
	}
}

// Enqueue adds a task to the default queue
func (p *RedisProducer) Enqueue(ctx context.Context, task *Task) (string, error) {
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now().Unix()
	task.Queue = p.getQueueName(task.Queue)
	task.TaskType = task.Handler

	if task.MaxRetry == 0 {
		task.MaxRetry = p.config.MaxRetry
	}

	payload, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("marshal task error: %w", err)
	}

	// Create asynq task using NewTask function
	t := asynq.NewTask(task.TaskType, payload)

	// Enqueue the task
	_, err = p.client.EnqueueContext(ctx, t, asynq.Queue(task.Queue), asynq.MaxRetry(task.MaxRetry))
	if err != nil {
		return "", fmt.Errorf("enqueue error: %w", err)
	}

	logx.Infof("Task %s enqueued to queue %s", task.ID, task.Queue)
	return task.ID, nil
}

// EnqueueWithPriority adds a task with specific priority
func (p *RedisProducer) EnqueueWithPriority(ctx context.Context, task *Task, priority Priority) (string, error) {
	task.Priority = priority
	return p.Enqueue(ctx, task)
}

// EnqueueDelayed adds a task that will be executed after the specified delay
func (p *RedisProducer) EnqueueDelayed(ctx context.Context, task *Task, delay time.Duration) (string, error) {
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now().Unix()
	task.Queue = p.getQueueName(task.Queue)
	task.TaskType = task.Handler
	task.Delay = int64(delay.Seconds())

	if task.MaxRetry == 0 {
		task.MaxRetry = p.config.MaxRetry
	}

	payload, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("marshal task error: %w", err)
	}

	// Create asynq task using NewTask function
	t := asynq.NewTask(task.TaskType, payload)

	// Enqueue the task with delay
	_, err = p.client.EnqueueContext(ctx, t, asynq.Queue(task.Queue), asynq.MaxRetry(task.MaxRetry), asynq.ProcessIn(delay))
	if err != nil {
		return "", fmt.Errorf("enqueue delayed error: %w", err)
	}

	logx.Infof("Task %s enqueued to queue %s with delay %v", task.ID, task.Queue, delay)
	return task.ID, nil
}

// EnqueueToQueue adds a task to a specific queue
func (p *RedisProducer) EnqueueToQueue(ctx context.Context, task *Task, queueName string) (string, error) {
	task.Queue = queueName
	return p.Enqueue(ctx, task)
}

// BatchEnqueue adds multiple tasks to the queue
func (p *RedisProducer) BatchEnqueue(ctx context.Context, tasks []*Task) ([]string, error) {
	taskIDs := make([]string, 0, len(tasks))
	var lastErr error

	for _, task := range tasks {
		id, err := p.Enqueue(ctx, task)
		if err != nil {
			lastErr = err
			logx.Errorf("BatchEnqueue error for task: %v", err)
			continue
		}
		taskIDs = append(taskIDs, id)
	}

	return taskIDs, lastErr
}

// CancelTask cancels a pending task
func (p *RedisProducer) CancelTask(ctx context.Context, taskID string) error {
	// Use Inspector to delete tasks
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:     p.config.RedisOpt.Addr,
		Password: p.config.RedisOpt.Password,
		DB:       p.config.RedisOpt.DB,
	})
	defer inspector.Close()

	// Try to delete from all queues
	queues := []string{"critical", "default", "low"}
	for _, queue := range queues {
		err := inspector.DeleteTask(taskID, queue)
		if err == nil {
			logx.Infof("Task %s cancelled from queue %s", taskID, queue)
			return nil
		}
	}

	return fmt.Errorf("task %s not found", taskID)
}

// getQueueName returns the queue name or default
func (p *RedisProducer) getQueueName(queue string) string {
	if queue != "" {
		return queue
	}
	return "default"
}

