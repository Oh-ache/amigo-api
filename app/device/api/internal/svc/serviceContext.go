package svc

import (
	"amigo-api/app/device/api/internal/config"
	"amigo-api/app/device/rpc/device"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	DeviceRpcClient   device.Device
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		DeviceRpcClient:   device.NewDevice(zrpc.MustNewClient(c.DeviceRpcConf)),
	}
}
