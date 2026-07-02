// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"amigo-api/app/job/mqueue/internal/svc"
	"amigo-api/app/job/mqueue/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnqueueDelayedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnqueueDelayedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnqueueDelayedLogic {
	return &EnqueueDelayedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnqueueDelayedLogic) EnqueueDelayed(req *types.EnqueueReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
