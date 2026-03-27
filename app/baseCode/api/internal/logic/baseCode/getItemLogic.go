package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetItemLogic {
	return &GetItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetItemLogic) GetItem(req *types.GetBaseCodeItemReq) (resp *types.GetBaseCodeItemResp, err error) {
	var pbReq pb.GetBaseCodeItemReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	pbResp, err := l.svcCtx.BaseCodeRpc.GetBaseCodeItem(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	var apiResp types.GetBaseCodeItemResp
	if err := copier.Copy(&apiResp, pbResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}
