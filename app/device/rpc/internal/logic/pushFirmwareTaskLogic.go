package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushFirmwareTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushFirmwareTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushFirmwareTaskLogic {
	return &PushFirmwareTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushFirmwareTaskLogic) PushFirmwareTask(in *pb.PushFirmwareTaskReq) (*pb.FirmwareTaskResp, error) {
	_, err := l.svcCtx.FirmwareModel.FindOne(l.ctx, in.FirmwareId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	var lastTask model.FirmwareTask
	for _, deviceId := range in.DeviceIds {
		task := model.FirmwareTask{
			FirmwareId: in.FirmwareId,
			DeviceId:   deviceId,
			Status:     1,
			IsDelete:   2,
		}
		if _, err := l.svcCtx.FirmwareTaskModel.Insert(l.ctx, &task); err != nil {
			return nil, err
		}
		lastTask = task
	}

	return &pb.FirmwareTaskResp{
		FirmwareTaskId: lastTask.FirmwareTaskId,
		FirmwareId:     lastTask.FirmwareId,
		DeviceId:       lastTask.DeviceId,
		Status:         lastTask.Status,
		Progress:       lastTask.Progress,
		ErrorMsg:       lastTask.ErrorMsg,
		StartedAt:      lastTask.StartedAt,
		CompletedAt:    lastTask.CompletedAt,
		CreateTime:     lastTask.CreateTime,
		UpdateTime:     lastTask.UpdateTime,
	}, nil
}
