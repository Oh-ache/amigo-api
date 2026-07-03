package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFirmwareTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFirmwareTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFirmwareTaskLogic {
	return &ListFirmwareTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFirmwareTaskLogic) ListFirmwareTask(in *pb.ListFirmwareTaskReq) (*pb.ListFirmwareTaskResp, error) {
	search := &model.FirmwareTaskSearch{
		FirmwareId: in.FirmwareId,
		DeviceId:   in.DeviceId,
		Status:     in.Status,
		Page:       in.Page,
		PageSize:   in.PageSize,
	}

	list, total, err := l.svcCtx.FirmwareTaskModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListFirmwareTaskResp{Total: total}
	for _, item := range list {
		resp.List = append(resp.List, &pb.FirmwareTaskResp{
			FirmwareTaskId: item.FirmwareTaskId,
			FirmwareId:     item.FirmwareId,
			DeviceId:       item.DeviceId,
			Status:         item.Status,
			Progress:       item.Progress,
			ErrorMsg:       item.ErrorMsg,
			StartedAt:      item.StartedAt,
			CompletedAt:    item.CompletedAt,
			CreateTime:     item.CreateTime,
			UpdateTime:     item.UpdateTime,
		})
	}

	return resp, nil
}
