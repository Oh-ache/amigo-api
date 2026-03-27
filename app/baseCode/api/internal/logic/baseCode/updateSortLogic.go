package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
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
	var pbReq pb.BaseCodeSortResp
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BaseCodeRpc.UpdateBaseCodeSort(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
