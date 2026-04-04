package queue

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	queueOnce      sync.Once
	globalProducer Producer
	globalConsumer Consumer
)

// InitGlobalQueue 初始化全局队列客户端
func InitGlobalQueue(redisOpt *redis.Options, config *QueueConfig) error {
	var err error
	queueOnce.Do(func() {
		rdb := redis.NewClient(redisOpt)
		queueClient := NewRedisQueueClient(rdb, config)
		globalProducer = NewRedisProducer(queueClient)
	})
	return err
}

// GetProducer 获取全局生产者
func GetProducer() Producer {
	return globalProducer
}

// GetConsumer 获取全局消费者
func GetConsumer() Consumer {
	return globalConsumer
}

