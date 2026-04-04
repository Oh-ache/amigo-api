package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"amigo-api/app/job/queue/internal/types"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// 消费者接口
type Consumer interface {
	RegisterHandler(name string, handler types.Handler)
	Start(ctx context.Context) error
	Stop() error
	ProcessOne(ctx context.Context, queue string) error
}

// 消费者实现
type RedisConsumer struct {
	client   *RedisQueueClient
	handlers map[string]types.Handler
	config   *ConsumerConfig
	stopChan chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

type ConsumerConfig struct {
	Concurrency  int
	PollInterval time.Duration
	WorkerPool   int
}

func NewRedisConsumer(client *RedisQueueClient, config *ConsumerConfig) *RedisConsumer {
	if config.Concurrency <= 0 {
		config.Concurrency = 10
	}
	if config.PollInterval <= 0 {
		config.PollInterval = time.Second
	}
	if config.WorkerPool <= 0 {
		config.WorkerPool = 100
	}

	return &RedisConsumer{
		client:   client,
		handlers: make(map[string]types.Handler),
		config:   config,
		stopChan: make(chan struct{}),
	}
}

// 注册处理器
func (c *RedisConsumer) RegisterHandler(name string, handler types.Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[name] = handler
}

// 启动消费者
func (c *RedisConsumer) Start(ctx context.Context) error {
	// 启动延迟任务处理器
	go c.processDelayTasks(ctx)

	// 启动处理中任务恢复器
	go c.recoverProcessingTasks(ctx)

	// 启动指定数量的消费者协程
	for i := 0; i < c.config.Concurrency; i++ {
		c.wg.Add(1)
		go c.worker(ctx, i)
	}

	logx.Infof("Queue consumer started with %d workers", c.config.Concurrency)
	return nil
}

// 停止消费者
func (c *RedisConsumer) Stop() error {
	close(c.stopChan)
	c.wg.Wait()
	logx.Info("Queue consumer stopped")
	return nil
}

// 工作协程
func (c *RedisConsumer) worker(ctx context.Context, workerID int) {
	defer c.wg.Done()

	logx.Infof("Worker %d started", workerID)

	for {
		select {
		case <-c.stopChan:
			logx.Infof("Worker %d stopped", workerID)
			return
		default:
			// 处理延迟队列
			c.processDelayedTasks(ctx)

			// 处理普通队列
			if err := c.ProcessOne(ctx, c.client.config.DefaultQueue); err != nil {
				logx.Errorf("Worker %d process error: %v", workerID, err)
				time.Sleep(time.Second)
			}

			time.Sleep(c.config.PollInterval)
		}
	}
}

// 处理单个任务
func (c *RedisConsumer) ProcessOne(ctx context.Context, queue string) error {
	queueKey := c.client.getQueueKey(queue)
	processingKey := c.client.getProcessingKey(queue)

	// 从队列中获取任务（优先级最高的任务）
	result, err := c.client.rdb.ZPopMax(ctx, queueKey, 1).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("pop from queue error: %w", err)
	}

	if len(result) == 0 {
		return nil // 队列为空
	}

	taskID := result[0].Member.(string)

	// 获取任务数据
	taskKey := c.client.getTaskKey(taskID)
	taskData, err := c.client.rdb.Get(ctx, taskKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			// 任务数据不存在，可能是被清理了
			return nil
		}
		return fmt.Errorf("get task error: %w", err)
	}

	var task types.Task
	if err := json.Unmarshal(taskData, &task); err != nil {
		// 数据格式错误，删除任务
		c.client.rdb.Del(ctx, taskKey)
		return fmt.Errorf("unmarshal task error: %w", err)
	}

	// 将任务加入处理中集合（设置超时时间）
	now := time.Now().Unix()
	processingScore := float64(now + task.Timeout)
	z := redis.Z{
		Score:  processingScore,
		Member: taskID,
	}
	if err := c.client.rdb.ZAdd(ctx, processingKey, z).Err(); err != nil {
		return fmt.Errorf("add to processing set error: %w", err)
	}

	// 异步处理任务
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.handleTask(ctx, &task)
	}()

	return nil
}

// 处理任务
func (c *RedisConsumer) handleTask(ctx context.Context, task *types.Task) {
	startTime := time.Now()

	// 查找处理器
	c.mu.RLock()
	handler, ok := c.handlers[task.Handler]
	c.mu.RUnlock()

	var err error
	if ok {
		err = handler.Handle(ctx, task)
	} else {
		err = fmt.Errorf("handler '%s' not found", task.Handler)
	}

	duration := time.Since(startTime).Milliseconds()

	// 更新任务状态
	if err != nil {
		logx.Errorf("Handle task %s error: %v", task.ID, err)
		c.handleFailedTask(ctx, task, err)
	} else {
		logx.Infof("Handle task %s success, duration: %dms", task.ID, duration)
		c.handleSuccessTask(ctx, task, duration)
	}
}

// 处理成功任务
func (c *RedisConsumer) handleSuccessTask(ctx context.Context, task *types.Task, duration int64) {
	// 从处理中集合移除
	processingKey := c.client.getProcessingKey(task.Queue)
	c.client.rdb.ZRem(ctx, processingKey, task.ID)

	// 删除任务数据
	taskKey := c.client.getTaskKey(task.ID)
	c.client.rdb.Del(ctx, taskKey)

	// 记录结果（可选）
	result := &types.TaskResult{
		TaskID:     task.ID,
		Status:     types.StatusCompleted,
		Duration:   duration,
		FinishedAt: time.Now().Unix(),
	}
	resultData, _ := json.Marshal(result)
	resultKey := fmt.Sprintf("%sresult:%s", c.client.config.Prefix, task.ID)
	c.client.rdb.SetEx(ctx, resultKey, string(resultData), 7*24*time.Hour)
}

// 处理失败任务
func (c *RedisConsumer) handleFailedTask(ctx context.Context, task *types.Task, err error) {
	// 从处理中集合移除
	processingKey := c.client.getProcessingKey(task.Queue)
	c.client.rdb.ZRem(ctx, processingKey, task.ID)

	task.RetryCount++

	// 检查重试次数
	if task.RetryCount >= task.MaxRetry {
		// 重试次数用尽，进入死信队列
		c.moveToDeadLetter(ctx, task, err)
	} else {
		// 重试
		c.retryTask(ctx, task, err)
	}
}

// 重试任务
func (c *RedisConsumer) retryTask(ctx context.Context, task *types.Task, err error) {
	// 延迟重试（指数退避）
	retryDelay := time.Duration(1<<uint(task.RetryCount)) * time.Second

	// 更新任务数据
	taskData, _ := json.Marshal(task)
	taskKey := c.client.getTaskKey(task.ID)
	c.client.rdb.SetEx(ctx, taskKey, string(taskData), 24*time.Hour)

	// 加入延迟队列
	delayKey := c.client.getDelaySetKey(c.client.config.RetryQueue)
	executeAt := time.Now().Add(retryDelay).Unix()

	z := redis.Z{
		Score:  float64(executeAt),
		Member: task.ID,
	}
	c.client.rdb.ZAdd(ctx, delayKey, z)

	logx.Infof("Task %s will retry after %v (retry %d/%d)",
		task.ID, retryDelay, task.RetryCount, task.MaxRetry)
}

// 转移到死信队列
func (c *RedisConsumer) moveToDeadLetter(ctx context.Context, task *types.Task, err error) {
	// 更新任务状态
	taskData, _ := json.Marshal(task)
	taskKey := c.client.getTaskKey(task.ID)
	c.client.rdb.SetEx(ctx, taskKey, string(taskData), 7*24*time.Hour)

	// 加入死信队列
	deadLetterKey := c.client.getQueueKey(c.client.config.DeadLetterQueue)
	z := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: task.ID,
	}
	c.client.rdb.ZAdd(ctx, deadLetterKey, z)

	// 记录失败结果
	result := &types.TaskResult{
		TaskID:     task.ID,
		Status:     types.StatusDeadLetter,
		Error:      err.Error(),
		FinishedAt: time.Now().Unix(),
	}
	resultData, _ := json.Marshal(result)
	resultKey := fmt.Sprintf("%sresult:%s", c.client.config.Prefix, task.ID)
	c.client.rdb.SetEx(ctx, resultKey, string(resultData), 30*24*time.Hour)

	logx.Errorf("Task %s moved to dead letter queue: %v", task.ID, err)
}

// 处理延迟任务
func (c *RedisConsumer) processDelayTasks(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			c.processDelayedTasks(ctx)
		}
	}
}

// 处理延迟任务
func (c *RedisConsumer) processDelayedTasks(ctx context.Context) {
	queues := []string{c.client.config.DefaultQueue, c.client.config.RetryQueue}

	for _, queue := range queues {
		c.processDelayedQueue(ctx, queue)
	}
}

// 处理指定队列的延迟任务
func (c *RedisConsumer) processDelayedQueue(ctx context.Context, queue string) {
	delayKey := c.client.getDelaySetKey(queue)
	now := time.Now().Unix()

	// 获取到期的延迟任务
	opt := &redis.ZRangeBy{
		Min:   "0",
		Max:   fmt.Sprintf("%d", now),
		Count: 10,
	}

	taskIDs, err := c.client.rdb.ZRangeByScore(ctx, delayKey, opt).Result()
	if err != nil {
		logx.Errorf("Get delayed tasks error: %v", err)
		return
	}

	if len(taskIDs) == 0 {
		return
	}

	// 从延迟集合中移除
	interfaceIDs := make([]interface{}, len(taskIDs))
	for i, id := range taskIDs {
		interfaceIDs[i] = id
	}
	c.client.rdb.ZRem(ctx, delayKey, interfaceIDs...)

	// 将任务加入执行队列
	for _, taskID := range taskIDs {
		// 获取任务数据
		taskKey := c.client.getTaskKey(taskID)
		taskData, err := c.client.rdb.Get(ctx, taskKey).Bytes()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			logx.Errorf("Get delayed task data error: %v", err)
			continue
		}

		var task types.Task
		if err := json.Unmarshal(taskData, &task); err != nil {
			logx.Errorf("Unmarshal delayed task error: %v", err)
			c.client.rdb.Del(ctx, taskKey)
			continue
		}

		// 加入执行队列
		queueKey := c.client.getQueueKey(queue)
		z := redis.Z{
			Score:  float64(task.Priority),
			Member: taskID,
		}
		if err := c.client.rdb.ZAdd(ctx, queueKey, z).Err(); err != nil {
			logx.Errorf("Add delayed task to queue error: %v", err)
		} else {
			logx.Infof("Delayed task %s moved to queue %s", taskID, queue)
		}
	}
}

// 恢复处理中超时的任务
func (c *RedisConsumer) recoverProcessingTasks(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			queues := []string{c.client.config.DefaultQueue, c.client.config.RetryQueue}

			for _, queue := range queues {
				c.recoverQueueProcessingTasks(ctx, queue)
			}
		}
	}
}

// 恢复指定队列的处理中任务
func (c *RedisConsumer) recoverQueueProcessingTasks(ctx context.Context, queue string) {
	processingKey := c.client.getProcessingKey(queue)
	now := time.Now().Unix()

	// 获取超时的处理中任务
	opt := &redis.ZRangeBy{
		Min:   "0",
		Max:   fmt.Sprintf("%d", now),
		Count: 100,
	}

	taskIDs, err := c.client.rdb.ZRangeByScore(ctx, processingKey, opt).Result()
	if err != nil {
		logx.Errorf("Get processing tasks error: %v", err)
		return
	}

	if len(taskIDs) == 0 {
		return
	}

	// 从处理中集合移除
	interfaceIDs := make([]interface{}, len(taskIDs))
	for i, id := range taskIDs {
		interfaceIDs[i] = id
	}
	c.client.rdb.ZRem(ctx, processingKey, interfaceIDs...)

	// 重新加入队列
	for _, taskID := range taskIDs {
		// 获取任务数据
		taskKey := c.client.getTaskKey(taskID)
		taskData, err := c.client.rdb.Get(ctx, taskKey).Bytes()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			logx.Errorf("Get processing task data error: %v", err)
			continue
		}

		var task types.Task
		if err := json.Unmarshal(taskData, &task); err != nil {
			logx.Errorf("Unmarshal processing task error: %v", err)
			c.client.rdb.Del(ctx, taskKey)
			continue
		}

		// 检查重试次数
		if task.RetryCount >= task.MaxRetry {
			c.moveToDeadLetter(ctx, &task, fmt.Errorf("task timeout in processing"))
		} else {
			// 重新入队
			queueKey := c.client.getQueueKey(queue)
			z := redis.Z{
				Score:  float64(task.Priority),
				Member: taskID,
			}
			c.client.rdb.ZAdd(ctx, queueKey, z)
			logx.Infof("Recovered timeout task %s to queue %s", taskID, queue)
		}
	}
}

