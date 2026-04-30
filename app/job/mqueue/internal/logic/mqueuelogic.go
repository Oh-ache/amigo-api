package logic

import (
	"context"

	"amigo-api/app/job/mqueue/internal/svc"
	"amigo-api/app/job/mqueue/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MqueueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMqueueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MqueueLogic {
	return &MqueueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MqueueLogic) Mqueue(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
