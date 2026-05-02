package logic

import (
	"context"

	"amigo-api/app/ai/model"
	"amigo-api/app/ai/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskLogic) GetTask(in *pb.GetTaskReq) (*pb.GetTaskResp, error) {
	task, err := l.svcCtx.AiTaskModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if task.UserId != in.UserId {
		return nil, model.ErrNotFound
	}

	return &pb.GetTaskResp{
		Id:           task.Id,
		UserId:       task.UserId,
		TaskId:       task.TaskId,
		TaskType:     stringToTaskType(task.TaskType),
		Prompt:       task.Prompt,
		RequestInfo:  task.RequestInfo,
		ResponseInfo: task.ResponseInfo,
		ResultUrl:    task.ResultUrl,
		Status:       int32(task.Status),
		ErrorMsg:     task.ErrorMsg,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
	}, nil
}

func stringToTaskType(s string) pb.TaskType {
	switch s {
	case "text_to_image":
		return pb.TaskType_TEXT_TO_IMAGE
	case "video":
		return pb.TaskType_VIDEO
	case "audio":
		return pb.TaskType_AUDIO
	default:
		return pb.TaskType_TEXT_TO_IMAGE
	}
}