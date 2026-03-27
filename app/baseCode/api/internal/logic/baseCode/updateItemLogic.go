package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateItemLogic {
	return &UpdateItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateItemLogic) UpdateItem(req *types.UpdateBaseCodeItemReq) (resp *types.EmptyResp, err error) {
	var pbReq pb.BaseCodeItemResp
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BaseCodeRpc.UpdateBaseCodeItem(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
