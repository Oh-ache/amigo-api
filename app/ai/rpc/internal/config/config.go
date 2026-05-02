package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource string
	}
	Cache []redis.RedisConf
	AiRedis struct {
		Host string
		Type string
		Pass string
		Key  string
	}
	BaseCodeRpcConf zrpc.RpcClientConf
	SdkRpcConf      zrpc.RpcClientConf
}