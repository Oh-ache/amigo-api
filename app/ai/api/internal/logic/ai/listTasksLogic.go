package ai

import (
	"context"
	"errors"

	"amigo-api/app/ai/api/internal/svc"
	"amigo-api/app/ai/api/internal/types"
	"amigo-api/app/ai/rpc/airpc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTasksLogic {
	return &ListTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTasksLogic) ListTasks(req *types.TaskListReq) (resp *types.TaskListResp, err error) {
	userId := getUserIdFromContext(l.ctx)
	if userId == 0 {
		return nil, errors.New("invalid user")
	}

	param := &airpc.ListTasksReq{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		TaskId:   req.TaskId,
		TaskType: pb.TaskType(stringToTaskTypeInt(req.TaskType)),
		Status:   int32(req.Status),
		UserId:   userId,
	}

	rpcResp, err := l.svcCtx.AiRpc.ListTasks(l.ctx, param)
	if err != nil {
		return nil, err
	}

	list := make([]types.GetTaskResp, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, types.GetTaskResp{
			Id:        item.Id,
			TaskId:    item.TaskId,
			TaskType:  taskTypeToStringInt(int32(item.TaskType)),
			Prompt:    item.Prompt,
			Status:    int(item.Status),
			ResultUrl: item.ResultUrl,
			ErrorMsg:  item.ErrorMsg,
			CreatedAt: item.CreatedAt,
		})
	}

	resp = &types.TaskListResp{
		Total: rpcResp.Total,
		List:  list,
	}

	return resp, nil
}

func stringToTaskTypeInt(s string) int32 {
	switch s {
	case "text_to_image":
		return 0
	case "video":
		return 1
	case "audio":
		return 2
	default:
		return 0
	}
}