package mqueue

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	mqueueOnce      sync.Once
	globalProducer  Producer
	globalConsumer  Consumer
	globalProducer4 *RedisProducer
)

// InitGlobalMQueue initializes the global message queue
func InitGlobalMQueue(redisOpt *redis.Options, config *QueueConfig) error {
	var err error
	mqueueOnce.Do(func() {
		// Create producer
		globalProducer4 = NewRedisProducer(redisOpt, config)
		globalProducer = globalProducer4

		// Create consumer
		globalConsumer = NewRedisConsumer(redisOpt, config)
	})
	return err
}

// InitGlobalMQueueWithDefaults initializes with default configuration
func InitGlobalMQueueWithDefaults(redisOpt *redis.Options) error {
	config := GetDefaultQueueConfig()
	config.RedisOpt = redisOpt
	return InitGlobalMQueue(redisOpt, config)
}

// GetProducer returns the global producer
func GetProducer() Producer {
	return globalProducer
}

// GetConsumer returns the global consumer
func GetConsumer() Consumer {
	return globalConsumer
}

// GetProducerWithContext returns producer with context support
func GetProducerWithContext(ctx context.Context) *RedisProducer {
	return globalProducer4
}

// Shutdown gracefully shuts down the message queue
func Shutdown() error {
	if globalConsumer != nil {
		return globalConsumer.Stop()
	}
	return nil
}

