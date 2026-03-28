package message

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckCodeLogic {
	return &CheckCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckCodeLogic) CheckCode(req *types.CheckCodeReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
