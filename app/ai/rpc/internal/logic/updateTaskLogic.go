package logic

import (
	"context"

	"amigo-api/app/ai/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTaskLogic) UpdateTask(in *pb.UpdateTaskReq) (*pb.UpdateTaskResp, error) {
	task, err := l.svcCtx.AiTaskModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if in.TaskId != "" {
		task.TaskId = in.TaskId
	}
	if in.ResponseInfo != "" {
		task.ResponseInfo = in.ResponseInfo
	}
	if in.ResultUrl != "" {
		task.ResultUrl = in.ResultUrl
	}
	if in.Status > 0 {
		task.Status = int(in.Status)
	}
	if in.ErrorMsg != "" {
		task.ErrorMsg = in.ErrorMsg
	}

	err = l.svcCtx.AiTaskModel.Update(l.ctx, task)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTaskResp{
		Success: true,
	}, nil
}