package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceLogic {
	return &GetDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceLogic) GetDevice(in *pb.GetDeviceReq) (*pb.DeviceResp, error) {
	var device *model.Device
	var err error

	if in.DeviceId != 0 {
		device, err = l.svcCtx.DeviceModel.FindOne(l.ctx, in.DeviceId)
	} else if in.MacAddress != "" {
		device, err = l.svcCtx.DeviceModel.FindOneByMacAddress(l.ctx, in.MacAddress)
	} else {
		return nil, model.ErrNotFound
	}

	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	var resp pb.DeviceResp
	if err := copier.Copy(&resp, device); err != nil {
		return nil, err
	}

	return &resp, nil
}
