package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSortLogic {
	return &UpdateSortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSortLogic) UpdateSort(req *types.UpdateBaseCodeSortReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
