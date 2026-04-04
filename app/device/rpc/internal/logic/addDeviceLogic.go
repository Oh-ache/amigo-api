package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDeviceLogic {
	return &AddDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddDeviceLogic) AddDevice(in *pb.AddDeviceReq) (*pb.DeviceResp, error) {
	var m model.Device
	if err := copier.Copy(&m, in); err != nil {
		return nil, err
	}

	if m.IsDelete == 0 {
		m.IsDelete = 2
	}

	isDuplicate, err := l.svcCtx.DeviceModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	result, err := l.svcCtx.DeviceModel.Insert(l.ctx, &m)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.DeviceId = uint64(id)

	var resp pb.DeviceResp
	if err := copier.Copy(&resp, &m); err != nil {
		return nil, err
	}

	return &resp, nil
}
