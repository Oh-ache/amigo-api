// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"amigo-api/app/job/mqueue/internal/svc"
	"amigo-api/app/job/mqueue/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnqueueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnqueueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnqueueLogic {
	return &EnqueueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnqueueLogic) Enqueue(req *types.EnqueueReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
