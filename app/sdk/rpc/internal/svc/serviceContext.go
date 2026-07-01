package svc

import (
	"context"
	"sync"

	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/app/sdk/rpc/internal/config"
	"amigo-api/common/pb"
	"amigo-api/common/utils/plug/objectsave/client"
	"amigo-api/common/utils/plug/objectsave/factory"
	"amigo-api/common/utils/plug/objectsave/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis
	BaseCodeRpc basecode.BaseCode
	OssClient   client.StorageClient
	initOnce    sync.Once
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		BaseCodeRpc: basecode.NewBaseCode(zrpc.MustNewClient(c.BaseCodeRpcConf)),
	}
}

func (s *ServiceContext) GetOssClient() (client.StorageClient, error) {
	var initErr error
	s.initOnce.Do(func() {
		f := factory.NewStorageFactory()
		storageConfig := &model.StorageConfig{
			Type: "oss",
			OssConfig: &model.OssConfig{
				Endpoint:        getBaseCodeFromRpc(s.BaseCodeRpc, "sdk", "ali.oss.endpoint"),
				AccessKeyId:     getBaseCodeFromRpc(s.BaseCodeRpc, "sdk", "ali.accessKey"),
				AccessKeySecret: getBaseCodeFromRpc(s.BaseCodeRpc, "sdk", "ali.accessKeySecret"),
				Bucket:          getBaseCodeFromRpc(s.BaseCodeRpc, "sdk", "ali.oss.bucket"),
				Region:          getBaseCodeFromRpc(s.BaseCodeRpc, "sdk", "ali.oss.region"),
			},
		}
		s.OssClient, initErr = f.CreateClient(storageConfig)
	})
	return s.OssClient, initErr
}

func getBaseCodeFromRpc(rpc basecode.BaseCode, sortKey, key string) string {
	item, _ := rpc.GetBaseCode(context.Background(), &pb.GetBaseCodeReq{
		SortKey: sortKey,
		Key:     key,
	})
	if item != nil {
		return item.Content
	}
	return ""
}
