package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFirmwareLogic {
	return &GetFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFirmwareLogic) GetFirmware(in *pb.GetFirmwareReq) (*pb.FirmwareResp, error) {
	data, err := l.svcCtx.FirmwareModel.FindOne(l.ctx, in.FirmwareId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
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
