package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"amigo-api/app/job/queue/internal/config"
	"amigo-api/app/job/queue/queue"
	"amigo-api/app/job/queue/queue/handler"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	configFile := flag.String("f", "etc/queue.yaml", "config file")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 设置日志
	logx.MustSetup(c.Log)

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

	// 创建消费者配置
	consumerConfig := &queue.ConsumerConfig{
		Concurrency:  c.Queue.Concurrency,
		PollInterval: time.Duration(c.Queue.PollInterval) * time.Second,
		WorkerPool:   100,
	}

	// 创建消费者
	consumer := queue.NewRedisConsumer(queueClient, consumerConfig)

	// 注册处理器
	registerHandlers(consumer)

	// 启动消费者
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		logx.Errorf("Failed to start consumer: %v", err)
		return
	}

	// 优雅关闭
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	fmt.Println("\nShutting down queue worker...")

	consumer.Stop()
	cancel()

	time.Sleep(2 * time.Second)
	fmt.Println("Queue worker stopped")
}

func registerHandlers(consumer queue.Consumer) {
	// 注册所有任务处理器
	consumer.RegisterHandler("send_sms", &handler.SendSmsHandler{})
}
