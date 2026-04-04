package logic

import (
	"context"

	"amigo-api/app/device/model"

	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDeviceLogic {
	return &UpdateDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDeviceLogic) UpdateDevice(in *pb.DeviceResp) (*pb.DeviceResp, error) {
	_, err := l.svcCtx.DeviceModel.FindOne(l.ctx, in.DeviceId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	var m model.Device
	if err := copier.Copy(&m, in); err != nil {
		return nil, err
	}

	isDuplicate, err := l.svcCtx.DeviceModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	if err := l.svcCtx.DeviceModel.Update(l.ctx, &m); err != nil {
		return nil, err
	}

	var resp pb.DeviceResp
	if err := copier.Copy(&resp, &m); err != nil {
		return nil, err
	}

	return &resp, nil
}
