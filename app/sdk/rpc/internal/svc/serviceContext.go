package svc

import (
	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/app/sdk/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis

	BaseCodeRpc basecode.BaseCode
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
