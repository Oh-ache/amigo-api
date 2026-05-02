package svc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"amigo-api/app/job/mqueue/internal/config"
	mqhandler "amigo-api/app/job/mqueue/internal/handler/mqueue"
	"amigo-api/app/ai/rpc/airpc"
	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/app/sdk/rpc/sdk"
	asynqmqueue "amigo-api/common/mqueue"
	"amigo-api/common/pb"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	Consumer *asynqmqueue.RedisConsumer
	Producer *asynqmqueue.RedisProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisOpt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.MQueue.RedisHost, c.MQueue.RedisPort),
		Password: c.MQueue.RedisPass,
		DB:       c.MQueue.RedisDB,
	}

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

	consumer := asynqmqueue.NewRedisConsumer(redisOpt, mqConfig)
	producer := asynqmqueue.NewRedisProducer(redisOpt, mqConfig)

	registerHandlers(consumer)

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

func NewServiceContextWithHandler(c config.Config, aiHandler *mqhandler.AiTaskHandler) *ServiceContext {
	redisOpt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.MQueue.RedisHost, c.MQueue.RedisPort),
		Password: c.MQueue.RedisPass,
		DB:       c.MQueue.RedisDB,
	}

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

	consumer := asynqmqueue.NewRedisConsumer(redisOpt, mqConfig)
	producer := asynqmqueue.NewRedisProducer(redisOpt, mqConfig)

	registerHandlersWithAi(consumer, aiHandler)

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

func registerHandlers(consumer *asynqmqueue.RedisConsumer) {
	consumer.RegisterHandler("send_sms", mqhandler.NewSendSmsHandler())
	logx.Info("MQueue handlers registered: send_sms")
}

func registerHandlersWithAi(consumer *asynqmqueue.RedisConsumer, aiHandler *mqhandler.AiTaskHandler) {
	consumer.RegisterHandler("send_sms", mqhandler.NewSendSmsHandler())
	consumer.RegisterHandler("ai_task", aiHandler)
	logx.Info("MQueue handlers registered: send_sms, ai_task")
}

type baseCodeRpcClient struct {
	cli basecode.BaseCode
}

func (c *baseCodeRpcClient) GetBaseCode(ctx context.Context, req *pb.GetBaseCodeReq) (*pb.BaseCodeResp, error) {
	return c.cli.GetBaseCode(ctx, req)
}

type aiRpcClient struct {
	aiCli  airpc.AiRpc
	sdkCli sdk.Sdk
}

func (c *aiRpcClient) UpdateTask(ctx context.Context, req *pb.UpdateTaskReq) (*pb.UpdateTaskResp, error) {
	return c.aiCli.UpdateTask(ctx, req)
}

func (c *aiRpcClient) UploadUrl(ctx context.Context, req *pb.UploadUrlReq) (*pb.UploadUrlResp, error) {
	return c.sdkCli.UploadUrl(ctx, req)
}

func NewBaseCodeRpcClient(cli zrpc.Client) mqhandler.BaseCodeRpcClient {
	return &baseCodeRpcClient{
		cli: basecode.NewBaseCode(cli),
	}
}

func NewAiRpcClient(cli zrpc.Client, sdkCli sdk.Sdk) mqhandler.AiRpcClient {
	return &aiRpcClient{
		aiCli:  airpc.NewAiRpc(cli),
		sdkCli: sdkCli,
	}
}

func NewSdkRpcClient(cli zrpc.Client) sdk.Sdk {
	return sdk.NewSdk(cli)
}

func (s *ServiceContext) EnqueueTask(ctx context.Context, handler string, data map[string]interface{}) (string, error) {
	task := &asynqmqueue.Task{
		Handler: handler,
		Data:    data,
		Queue:   getQueueFromHandler(handler),
	}
	return s.Producer.Enqueue(ctx, task)
}

func (s *ServiceContext) EnqueueDelayedTask(ctx context.Context, handler string, data map[string]interface{}, delay string) (string, error) {
	task := &asynqmqueue.Task{
		Handler: handler,
		Data:    data,
		Queue:   getQueueFromHandler(handler),
	}

	d, err := time.ParseDuration(delay)
	if err != nil {
		return "", err
	}

	return s.Producer.EnqueueDelayed(ctx, task, d)
}

func getQueueFromHandler(handler string) string {
	switch strings.ToLower(handler) {
	case "email", "payment", "ai_task":
		return "critical"
	case "notification", "task":
		return "default"
	case "log", "analytics":
		return "low"
	default:
		return "default"
	}
}