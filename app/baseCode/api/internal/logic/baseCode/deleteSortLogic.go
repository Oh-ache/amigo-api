package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSortLogic {
	return &DeleteSortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSortLogic) DeleteSort(req *types.DeleteBaseCodeSortReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
