package svc

import (
	"amigo-api/app/ai/api/internal/config"
	"amigo-api/app/ai/rpc/airpc"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	AiRpc  airpc.AiRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		AiRpc:  airpc.NewAiRpc(zrpc.MustNewClient(c.AiRpcConf)),
	}
}