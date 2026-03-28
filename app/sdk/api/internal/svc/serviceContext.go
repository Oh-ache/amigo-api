package svc

import (
	"amigo-api/app/sdk/api/internal/config"
	"amigo-api/app/sdk/rpc/sdk"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	SdkRpcClient sdk.Sdk
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		SdkRpcClient: sdk.NewSdk(zrpc.MustNewClient(c.SdkRpcConf)),
	}
}
