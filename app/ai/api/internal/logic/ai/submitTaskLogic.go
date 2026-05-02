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

type SubmitTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTaskLogic {
	return &SubmitTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitTaskLogic) SubmitTask(req *types.SubmitTaskReq) (resp *types.SubmitTaskResp, err error) {
	userId := getUserIdFromContext(l.ctx)
	if userId == 0 {
		return nil, errors.New("invalid user")
	}

	param := &airpc.SubmitTaskReq{
		TaskType: stringToTaskType(req.TaskType),
		Prompt:   req.Prompt,
		Params:   req.Params,
		UserId:   userId,
	}

	rpcResp, err := l.svcCtx.AiRpc.SubmitTask(l.ctx, param)
	if err != nil {
		return nil, err
	}

	resp = &types.SubmitTaskResp{
		Id:     rpcResp.Id,
		TaskId: rpcResp.TaskId,
	}

	return resp, nil
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

func taskTypeToStringInt(t int32) string {
	switch t {
	case 0:
		return "text_to_image"
	case 1:
		return "video"
	case 2:
		return "audio"
	default:
		return "text_to_image"
	}
}

func getUserIdFromContext(ctx context.Context) int64 {
	return 0
}