package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
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
	var pbReq pb.ListBaseCodeSortReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	pbResp, err := l.svcCtx.BaseCodeRpc.ListBaseCodeSort(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	var apiResp types.ListBaseCodeSortResp
	if err := copier.Copy(&apiResp, pbResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}
