package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFirmwareLogic {
	return &UpdateFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFirmwareLogic) UpdateFirmware(in *pb.FirmwareResp) (*pb.FirmwareResp, error) {
	data, err := l.svcCtx.FirmwareModel.FindOne(l.ctx, in.FirmwareId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	data.Name = in.Name
	data.Changelog = in.Changelog
	data.IsForce = in.IsForce

	if err := l.svcCtx.FirmwareModel.Update(l.ctx, data); err != nil {
		return nil, err
	}

	return &pb.FirmwareResp{
		FirmwareId: data.FirmwareId,
		Name:       data.Name,
		Version:    data.Version,
		DeviceType: data.DeviceType,
		FileUrl:    data.FileUrl,
		FileSize:   data.FileSize,
		Md5:        data.Md5,
		Changelog:  data.Changelog,
		IsForce:    data.IsForce,
		IsDelete:   data.IsDelete,
		CreateTime: data.CreateTime,
		UpdateTime: data.UpdateTime,
	}, nil
}
