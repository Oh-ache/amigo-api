package svc

import (
	"amigo-api/app/baseCode/api/internal/config"
	"amigo-api/app/baseCode/rpc/basecode"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	BaseCodeRpcClient basecode.BaseCode
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		BaseCodeRpcClient: basecode.NewBaseCode(zrpc.MustNewClient(c.BaseCodeRpcConf)),
	}
}
