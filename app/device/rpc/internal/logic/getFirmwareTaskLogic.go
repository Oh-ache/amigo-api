package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFirmwareTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFirmwareTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFirmwareTaskLogic {
	return &GetFirmwareTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFirmwareTaskLogic) GetFirmwareTask(in *pb.GetFirmwareTaskReq) (*pb.FirmwareTaskResp, error) {
	data, err := l.svcCtx.FirmwareTaskModel.FindOne(l.ctx, in.FirmwareTaskId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &pb.FirmwareTaskResp{
		FirmwareTaskId: data.FirmwareTaskId,
		FirmwareId:     data.FirmwareId,
		DeviceId:       data.DeviceId,
		Status:         data.Status,
		Progress:       data.Progress,
		ErrorMsg:       data.ErrorMsg,
		StartedAt:      data.StartedAt,
		CompletedAt:    data.CompletedAt,
		CreateTime:     data.CreateTime,
		UpdateTime:     data.UpdateTime,
	}, nil
}
