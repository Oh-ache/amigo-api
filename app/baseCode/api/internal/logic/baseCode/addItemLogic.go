package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddItemLogic {
	return &AddItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddItemLogic) AddItem(req *types.AddBaseCodeItemReq) (resp *types.EmptyResp, err error) {
	var pbReq pb.AddBaseCodeItemReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BaseCodeRpc.AddBaseCodeItem(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
