package svc

import (
	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/app/sdk/api/internal/config"
	"amigo-api/app/sdk/api/internal/svc/baidutts"
	"amigo-api/app/sdk/rpc/sdk"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	SdkRpcClient  sdk.Sdk
	BaseCodeRpc   basecode.BaseCode
	BaiduTTS      *baidutts.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	baseCodeRpc := basecode.NewBaseCode(zrpc.MustNewClient(c.BaseCodeRpcConf))
	return &ServiceContext{
		Config:       c,
		SdkRpcClient: sdk.NewSdk(zrpc.MustNewClient(c.SdkRpcConf)),
		BaseCodeRpc:  baseCodeRpc,
		BaiduTTS:     baidutts.New(baseCodeRpc),
	}
}
