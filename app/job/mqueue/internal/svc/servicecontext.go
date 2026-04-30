package svc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"amigo-api/app/job/mqueue/internal/config"
	"amigo-api/app/job/mqueue/internal/handler/mqueue"
	asynqmqueue "amigo-api/common/mqueue"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config   config.Config
	Consumer *asynqmqueue.RedisConsumer
	Producer *asynqmqueue.RedisProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Create Redis options for mqueue
	redisOpt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.MQueue.RedisHost, c.MQueue.RedisPort),
		Password: c.MQueue.RedisPass,
		DB:       c.MQueue.RedisDB,
	}

	// Create mqueue config
	mqConfig := &asynqmqueue.QueueConfig{
		RedisOpt:      redisOpt,
		Queues:        c.MQueue.Queues,
		Concurrency:   c.MQueue.Concurrency,
		StrictMode:    false,
		SyncTimeout:   c.MQueue.SyncTimeout,
		RetryDelay:    c.MQueue.RetryDelay,
		MaxRetry:      c.MQueue.MaxRetry,
		Timeout:       c.MQueue.Timeout,
		DeadQueueName: c.MQueue.DeadQueue,
		ServerName:    c.MQueue.ServerName,
	}

	// Create consumer and producer
	consumer := asynqmqueue.NewRedisConsumer(redisOpt, mqConfig)
	producer := asynqmqueue.NewRedisProducer(redisOpt, mqConfig)

	// Register handlers
	registerHandlers(consumer)

	// Start the consumer
	ctx := context.Background()
	if err := consumer.Start(ctx); err != nil {
		logx.Errorf("Failed to start mqueue consumer: %v", err)
	} else {
		logx.Info("MQueue consumer started successfully")
	}

	return &ServiceContext{
		Config:   c,
		Consumer: consumer,
		Producer: producer,
	}
}

// registerHandlers registers all task handlers
func registerHandlers(consumer *asynqmqueue.RedisConsumer) {
	// Register sample handler
	consumer.RegisterHandler("sample", mqueue.NewSampleHandler())

	// Register email handler
	consumer.RegisterHandler("email", mqueue.NewEmailHandler())

	// Register notification handler
	consumer.RegisterHandler("notification", mqueue.NewNotificationHandler())

	// Register task handler
	consumer.RegisterHandler("task", mqueue.NewTaskHandler())

	logx.Info("MQueue handlers registered: sample, email, notification, task")
}

// EnqueueTask is a helper method to enqueue tasks
func (s *ServiceContext) EnqueueTask(ctx context.Context, handler string, data map[string]interface{}) (string, error) {
	task := &asynqmqueue.Task{
		Handler: handler,
		Data:    data,
		Queue:   getQueueFromHandler(handler),
	}
	return s.Producer.Enqueue(ctx, task)
}

// EnqueueDelayedTask is a helper method to enqueue delayed tasks
func (s *ServiceContext) EnqueueDelayedTask(ctx context.Context, handler string, data map[string]interface{}, delay string) (string, error) {
	task := &asynqmqueue.Task{
		Handler: handler,
		Data:    data,
		Queue:   getQueueFromHandler(handler),
	}

	// Parse duration string like "1m", "30s", "1h"
	d, err := time.ParseDuration(delay)
	if err != nil {
		return "", err
	}

	return s.Producer.EnqueueDelayed(ctx, task, d)
}

// getQueueFromHandler returns the appropriate queue based on handler name
func getQueueFromHandler(handler string) string {
	switch strings.ToLower(handler) {
	case "email", "payment":
		return "critical"
	case "notification", "task":
		return "default"
	case "log", "analytics":
		return "low"
	default:
		return "default"
	}
}

