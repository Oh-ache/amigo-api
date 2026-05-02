package ai

import (
	"context"
	"errors"

	"amigo-api/app/ai/api/internal/svc"
	"amigo-api/app/ai/api/internal/types"
	"amigo-api/app/ai/rpc/airpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskLogic) GetTask(req *types.GetTaskReq) (resp *types.GetTaskResp, err error) {
	userId := getUserIdFromContext(l.ctx)
	if userId == 0 {
		return nil, errors.New("invalid user")
	}

	param := &airpc.GetTaskReq{
		Id:     req.Id,
		UserId: userId,
	}

	rpcResp, err := l.svcCtx.AiRpc.GetTask(l.ctx, param)
	if err != nil {
		return nil, err
	}

	resp = &types.GetTaskResp{
		Id:        rpcResp.Id,
		TaskId:    rpcResp.TaskId,
		TaskType:  taskTypeToStringInt(int32(rpcResp.TaskType)),
		Prompt:    rpcResp.Prompt,
		Status:    int(rpcResp.Status),
		ResultUrl: rpcResp.ResultUrl,
		ErrorMsg:  rpcResp.ErrorMsg,
		CreatedAt: rpcResp.CreatedAt,
	}

	return resp, nil
}