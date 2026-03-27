package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSortLogic {
	return &ListSortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSortLogic) ListSort(req *types.ListBaseCodeSortReq) (resp *types.ListBaseCodeSortResp, err error) {
	// todo: add your logic here and delete this line

	return
}
