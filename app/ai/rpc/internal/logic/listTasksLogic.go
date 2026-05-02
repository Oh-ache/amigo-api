package logic

import (
	"context"

	"amigo-api/app/ai/model"
	"amigo-api/app/ai/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTasksLogic {
	return &ListTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTasksLogic) ListTasks(in *pb.ListTasksReq) (*pb.ListTasksResp, error) {
	search := &model.AiTaskSearch{
		TaskId:   in.TaskId,
		TaskType: taskTypeToStringInt(in.TaskType),
		Status:   int(in.Status),
		UserId:   in.UserId,
		Page:     int64(in.Page),
		PageSize: int64(in.PageSize),
	}

	if search.Page <= 0 {
		search.Page = 1
	}
	if search.PageSize <= 0 {
		search.PageSize = 10
	}

	tasks, total, err := l.svcCtx.AiTaskModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	list := make([]*pb.GetTaskResp, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, &pb.GetTaskResp{
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
		})
	}

	return &pb.ListTasksResp{
		Total: total,
		List:  list,
	}, nil
}

func taskTypeToStringInt(t pb.TaskType) string {
	switch t {
	case pb.TaskType_TEXT_TO_IMAGE:
		return "text_to_image"
	case pb.TaskType_VIDEO:
		return "video"
	case pb.TaskType_AUDIO:
		return "audio"
	default:
		return ""
	}
}