package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDeviceLogic {
	return &DeleteDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDeviceLogic) DeleteDevice(in *pb.DeleteDeviceReq) (*pb.DeleteDeviceResp, error) {
	var device *model.Device
	var err error

	if in.DeviceId != 0 {
		device, err = l.svcCtx.DeviceModel.FindOne(l.ctx, in.DeviceId)
	} else if in.MacAddress != "" {
		device, err = l.svcCtx.DeviceModel.FindOneByMacAddress(l.ctx, in.MacAddress)
	} else {
		return &pb.DeleteDeviceResp{Success: false}, nil
	}

	if err != nil {
		if err == model.ErrNotFound {
			return &pb.DeleteDeviceResp{Success: false}, nil
		}
		return nil, err
	}

	if err := l.svcCtx.DeviceModel.Delete(l.ctx, device.DeviceId); err != nil {
		return nil, err
	}

	return &pb.DeleteDeviceResp{Success: true}, nil
}
