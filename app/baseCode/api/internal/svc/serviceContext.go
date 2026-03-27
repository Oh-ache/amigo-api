package svc

import (
	"amigo-api/app/baseCode/api/internal/config"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	BaseCodeRpc  pb.BaseCodeClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.BaseCodeRpc)
	return &ServiceContext{
		Config:       c,
		BaseCodeRpc:  pb.NewBaseCodeClient(client.Conn()),
	}
}
