package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFirmwareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFirmwareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFirmwareLogic {
	return &ListFirmwareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFirmwareLogic) ListFirmware(in *pb.ListFirmwareReq) (*pb.ListFirmwareResp, error) {
	search := &model.FirmwareSearch{
		DeviceType: in.DeviceType,
		Version:    in.Version,
		Name:       in.Name,
		IsForce:    in.IsForce,
		IsDelete:   in.IsDelete,
		Page:       in.Page,
		PageSize:   in.PageSize,
	}

	list, total, err := l.svcCtx.FirmwareModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListFirmwareResp{Total: total}
	for _, item := range list {
		resp.List = append(resp.List, &pb.FirmwareResp{
			FirmwareId: item.FirmwareId,
			Name:       item.Name,
			Version:    item.Version,
			DeviceType: item.DeviceType,
			FileUrl:    item.FileUrl,
			FileSize:   item.FileSize,
			Md5:        item.Md5,
			Changelog:  item.Changelog,
			IsForce:    item.IsForce,
			IsDelete:   item.IsDelete,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}

	return resp, nil
}
