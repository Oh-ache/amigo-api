package mqueue

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// QueueConfig holds configuration for the message queue
type QueueConfig struct {
	RedisOpt      *redis.Options   // Redis connection options
	Queues        map[string]int   // Queue priorities (higher = more workers)
	Concurrency   int              // Number of concurrent workers
	StrictMode    bool             // Whether to strictly enforce queue priorities
	SyncTimeout   time.Duration    // Timeout for sync operations
	RetryDelay    time.Duration    // Base delay for retries (with exponential backoff)
	MaxRetry      int              // Maximum retry attempts
	Timeout       time.Duration    // Task execution timeout
	DeadQueueName string           // Name of the dead letter queue
	ServerName    string           // Server identifier for asynq
}

// AsynqConfig converts QueueConfig to asynq.Config
func (c *QueueConfig) AsynqConfig() *AsynqConfig {
	return &AsynqConfig{
		Queues:        c.Queues,
		Concurrency:   c.Concurrency,
		StrictMode:    c.StrictMode,
		SyncTimeout:   c.SyncTimeout,
		RetryDelay:    c.RetryDelay,
		MaxRetry:      c.MaxRetry,
		Timeout:       c.Timeout,
		DeadQueueName: c.DeadQueueName,
		ServerName:    c.ServerName,
	}
}

// AsynqConfig is the internal config for asynq
type AsynqConfig struct {
	Queues        map[string]int
	Concurrency   int
	StrictMode    bool
	SyncTimeout   time.Duration
	RetryDelay    time.Duration
	MaxRetry      int
	Timeout       time.Duration
	DeadQueueName string
	ServerName    string
}

