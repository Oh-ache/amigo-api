package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB    struct {
		DataSource string
	}
	Cache cache.CacheConf
	EMQX  EMQXConf
}

type EMQXConf struct {
	Broker         string `json:",optional"`
	ApiEndpoint    string `json:",optional"`
	Username       string `json:",optional"`
	Password       string `json:",optional"`
	ClientIdPrefix string `json:",default=amigo-device"`
	KeepAlive      int64  `json:",default=60"`
	ConnectTimeout int64  `json:",default=10"`
	AutoReconnect  bool   `json:",default=true"`
	CleanSession   bool   `json:",default=true"`
}
