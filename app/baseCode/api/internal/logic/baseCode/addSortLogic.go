package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSortLogic {
	return &AddSortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSortLogic) AddSort(req *types.AddBaseCodeSortReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
