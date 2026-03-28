package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource string
	}
	Cache           cache.CacheConf
	BaseCodeRpcConf zrpc.RpcClientConf
	JwtAuth         struct {
		AccessSecret string
		AccessExpire int64
	}
}
