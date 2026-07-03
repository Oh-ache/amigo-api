package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFirmwareLogic {
	return &AddFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddFirmwareLogic) AddFirmware(in *pb.AddFirmwareReq) (*pb.FirmwareResp, error) {
	var m model.Firmware
	m.Name = in.Name
	m.Version = in.Version
	m.DeviceType = in.DeviceType
	m.FileUrl = in.FileUrl
	m.FileSize = in.FileSize
	m.Md5 = in.Md5
	m.Changelog = in.Changelog
	m.IsForce = in.IsForce
	m.IsDelete = 2

	isDuplicate, err := l.svcCtx.FirmwareModel.CheckDuplicate(l.ctx, &m)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, model.ErrDuplicate
	}

	if _, err := l.svcCtx.FirmwareModel.Insert(l.ctx, &m); err != nil {
		return nil, err
	}

	return &pb.FirmwareResp{
		FirmwareId: m.FirmwareId,
		Name:       m.Name,
		Version:    m.Version,
		DeviceType: m.DeviceType,
		FileUrl:    m.FileUrl,
		FileSize:   m.FileSize,
		Md5:        m.Md5,
		Changelog:  m.Changelog,
		IsForce:    m.IsForce,
		IsDelete:   m.IsDelete,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}, nil
}
