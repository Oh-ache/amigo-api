package queue

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisQueueClient struct {
	rdb    *redis.Client
	config *QueueConfig
}

type QueueConfig struct {
	Prefix          string
	DefaultQueue    string
	RetryQueue      string
	DelayQueue      string
	DeadLetterQueue string
	MaxRetry        int
}

func NewRedisQueueClient(rdb *redis.Client, config *QueueConfig) *RedisQueueClient {
	return &RedisQueueClient{
		rdb:    rdb,
		config: config,
	}
}

// 生成队列键名
func (c *RedisQueueClient) getQueueKey(queue string) string {
	return fmt.Sprintf("%s%s", c.config.Prefix, queue)
}

// 生成任务键名
func (c *RedisQueueClient) getTaskKey(taskID string) string {
	return fmt.Sprintf("%stask:%s", c.config.Prefix, taskID)
}

// 生成延迟队列键名
func (c *RedisQueueClient) getDelaySetKey(queue string) string {
	return fmt.Sprintf("%sdelay:%s", c.config.Prefix, queue)
}

// 生成处理中队列键名
func (c *RedisQueueClient) getProcessingKey(queue string) string {
	return fmt.Sprintf("%sprocessing:%s", c.config.Prefix, queue)
}

