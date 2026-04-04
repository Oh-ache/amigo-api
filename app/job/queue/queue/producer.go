package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"amigo-api/app/job/queue/internal/types"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// 生产者接口
type Producer interface {
	Enqueue(ctx context.Context, task *types.Task) (string, error)
	EnqueueWithPriority(ctx context.Context, task *types.Task, priority types.Priority) (string, error)
	EnqueueDelayed(ctx context.Context, task *types.Task, delay time.Duration) (string, error)
	BatchEnqueue(ctx context.Context, tasks []*types.Task) ([]string, error)
	CancelTask(ctx context.Context, taskID string) error
}

// 生产者实现
type RedisProducer struct {
	client *RedisQueueClient
}

func NewRedisProducer(client *RedisQueueClient) Producer {
	return &RedisProducer{client: client}
}

// 入队任务
func (p *RedisProducer) Enqueue(ctx context.Context, task *types.Task) (string, error) {
	return p.enqueue(ctx, task, task.Priority)
}

// 带优先级入队
func (p *RedisProducer) EnqueueWithPriority(ctx context.Context, task *types.Task, priority types.Priority) (string, error) {
	task.Priority = priority
	return p.enqueue(ctx, task, priority)
}

// 延迟任务
func (p *RedisProducer) EnqueueDelayed(ctx context.Context, task *types.Task, delay time.Duration) (string, error) {
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now().Unix()
	task.Delay = int64(delay.Seconds())

	if task.Queue == "" {
		task.Queue = p.client.config.DefaultQueue
	}
	if task.MaxRetry == 0 {
		task.MaxRetry = p.client.config.MaxRetry
	}

	// 序列化任务
	taskData, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("marshal task error: %w", err)
	}

	// 存储任务数据
	taskKey := p.client.getTaskKey(task.ID)
	if err := p.client.rdb.Set(ctx, taskKey, taskData, 24*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("save task error: %w", err)
	}

	// 加入延迟集合
	delayKey := p.client.getDelaySetKey(task.Queue)
	executeAt := time.Now().Add(delay).Unix()

	z := redis.Z{
		Score:  float64(executeAt),
		Member: task.ID,
	}

	if err := p.client.rdb.ZAdd(ctx, delayKey, z).Err(); err != nil {
		p.client.rdb.Del(ctx, taskKey)
		return "", fmt.Errorf("add to delay set error: %w", err)
	}

	return task.ID, nil
}

// 批量入队
func (p *RedisProducer) BatchEnqueue(ctx context.Context, tasks []*types.Task) ([]string, error) {
	pipe := p.client.rdb.Pipeline()
	taskIDs := make([]string, 0, len(tasks))

	for _, task := range tasks {
		task.ID = uuid.New().String()
		task.CreatedAt = time.Now().Unix()

		if task.Queue == "" {
			task.Queue = p.client.config.DefaultQueue
		}
		if task.MaxRetry == 0 {
			task.MaxRetry = p.client.config.MaxRetry
		}

		// 序列化任务
		taskData, err := json.Marshal(task)
		if err != nil {
			return nil, err
		}

		// 存储任务
		taskKey := p.client.getTaskKey(task.ID)
		pipe.Set(ctx, taskKey, taskData, 24*time.Hour)

		// 加入队列
		queueKey := p.client.getQueueKey(task.Queue)
		z := redis.Z{
			Score:  float64(task.Priority),
			Member: task.ID,
		}
		pipe.ZAdd(ctx, queueKey, z)

		taskIDs = append(taskIDs, task.ID)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return taskIDs, nil
}

// 取消任务
func (p *RedisProducer) CancelTask(ctx context.Context, taskID string) error {
	// 获取任务信息
	taskKey := p.client.getTaskKey(taskID)
	taskData, err := p.client.rdb.Get(ctx, taskKey).Bytes()
	if err != nil {
		return err
	}

	var task types.Task
	if err := json.Unmarshal(taskData, &task); err != nil {
		return err
	}

	// 从队列中移除
	queueKey := p.client.getQueueKey(task.Queue)
	p.client.rdb.ZRem(ctx, queueKey, taskID)

	// 从延迟集合中移除
	delayKey := p.client.getDelaySetKey(task.Queue)
	p.client.rdb.ZRem(ctx, delayKey, taskID)

	// 从处理中队列移除
	processingKey := p.client.getProcessingKey(task.Queue)
	p.client.rdb.ZRem(ctx, processingKey, taskID)

	// 删除任务数据
	p.client.rdb.Del(ctx, taskKey)

	return nil
}

// 私有方法：入队
func (p *RedisProducer) enqueue(ctx context.Context, task *types.Task, priority types.Priority) (string, error) {
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now().Unix()

	if task.Queue == "" {
		task.Queue = p.client.config.DefaultQueue
	}
	if task.MaxRetry == 0 {
		task.MaxRetry = p.client.config.MaxRetry
	}

	// 序列化任务
	taskData, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("marshal task error: %w", err)
	}

	// 存储任务数据
	taskKey := p.client.getTaskKey(task.ID)
	if err := p.client.rdb.Set(ctx, taskKey, taskData, 24*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("save task error: %w", err)
	}

	// 加入队列（使用有序集合，score为优先级）
	queueKey := p.client.getQueueKey(task.Queue)
	z := redis.Z{
		Score:  float64(priority),
		Member: task.ID,
	}

	if err := p.client.rdb.ZAdd(ctx, queueKey, z).Err(); err != nil {
		p.client.rdb.Del(ctx, taskKey)
		return "", fmt.Errorf("add to queue error: %w", err)
	}

	return task.ID, nil
}

