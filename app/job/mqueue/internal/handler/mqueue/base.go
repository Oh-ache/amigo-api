package mqueue

import (
	"context"

	"amigo-api/common/mqueue"
	"amigo-api/common/pb"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(client *redis.Client) {
	RedisClient = client
}

type BaseCodeRpcClient interface {
	GetBaseCode(ctx context.Context, req *pb.GetBaseCodeReq) (*pb.BaseCodeResp, error)
}

type AiRpcClient interface {
	UpdateTask(ctx context.Context, req *pb.UpdateTaskReq) (*pb.UpdateTaskResp, error)
	UploadUrl(ctx context.Context, req *pb.UploadUrlReq) (*pb.UploadUrlResp, error)
}

type TaskHandler interface {
	Name() string
	Handle(ctx context.Context, task *mqueue.Task) error
}
