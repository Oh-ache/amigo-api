package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSortLogic {
	return &GetSortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSortLogic) GetSort(req *types.GetBaseCodeSortReq) (resp *types.GetBaseCodeSortResp, err error) {
	// todo: add your logic here and delete this line

	return
}
