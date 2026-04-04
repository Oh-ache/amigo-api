package svc

import (
	"time"

	"amigo-api/app/job/queue/internal/config"
	"amigo-api/app/job/queue/queue/handler"
	"amigo-api/common/queue"

	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config   config.Config
	Redis    *redis.Client
	Queue    queue.Producer
	Consumer queue.Consumer
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化Redis
	// rdb := redis.MustNewRedis(c.Redis.RedisConf)

	// 创建Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       c.Redis.DB,
	})

	// 队列配置
	queueConfig := &queue.QueueConfig{
		Prefix:          c.Queue.Prefix,
		DefaultQueue:    c.Queue.DefaultQueue,
		RetryQueue:      c.Queue.RetryQueue,
		DelayQueue:      c.Queue.DelayQueue,
		DeadLetterQueue: c.Queue.DeadLetterQueue,
		MaxRetry:        c.Queue.MaxRetry,
	}

	// 创建队列客户端
	queueClient := queue.NewRedisQueueClient(redisClient, queueConfig)

	// 创建生产者
	producer := queue.NewRedisProducer(queueClient)

	// 创建消费者配置
	consumerConfig := &queue.ConsumerConfig{
		Concurrency:  c.Queue.Concurrency,
		PollInterval: time.Duration(c.Queue.PollInterval) * time.Second,
		WorkerPool:   100,
	}

	// 创建消费者
	consumer := queue.NewRedisConsumer(queueClient, consumerConfig)

	// 初始化处理器的 Redis 客户端
	handler.InitRedis(redisClient)

	// 注册处理器
	consumer.RegisterHandler("send_sms", &handler.SendSmsHandler{})

	return &ServiceContext{
		Config:   c,
		Redis:    redisClient,
		Queue:    producer,
		Consumer: consumer,
	}
}
