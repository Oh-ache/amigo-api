package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddDeviceEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddDeviceEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDeviceEventLogic {
	return &AddDeviceEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddDeviceEventLogic) AddDeviceEvent(in *pb.AddDeviceEventReq) (*pb.DeviceEventResp, error) {
	m := &model.DeviceEvent{
		DeviceId:   in.DeviceId,
		EventType:  in.EventType,
		EventLevel: in.EventLevel,
		Title:      in.Title,
	}
	if in.Description != "" {
		m.Description.String = in.Description
		m.Description.Valid = true
	}
	if in.ExtraData != "" {
		m.ExtraData.String = in.ExtraData
		m.ExtraData.Valid = true
	}
	if in.Source != "" {
		m.Source = in.Source
	} else {
		m.Source = "device"
	}
	if m.EventLevel == "" {
		m.EventLevel = "info"
	}
	if m.IsDelete == 0 {
		m.IsDelete = 2
	}

	result, err := l.svcCtx.DeviceEventModel.Insert(l.ctx, m)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.DeviceEventId = uint64(id)

	return toPbDeviceEvent(m), nil
}
