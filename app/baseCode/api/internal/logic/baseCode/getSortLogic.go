package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
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
	var pbReq pb.GetBaseCodeSortReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	pbResp, err := l.svcCtx.BaseCodeRpc.GetBaseCodeSort(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	var apiResp types.GetBaseCodeSortResp
	if err := copier.Copy(&apiResp, pbResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}
