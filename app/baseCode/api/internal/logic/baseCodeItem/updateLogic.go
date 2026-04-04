package baseCodeItem

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

)
type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateBaseCodeItemReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	// First, get the existing item to populate the BaseCodeItemResp
	getReq := &pb.GetBaseCodeItemReq{
		BaseCodeItemId: req.BaseCodeItemId,
	}
	existingItem, err := l.svcCtx.BaseCodeRpcClient.GetBaseCodeItem(l.ctx, getReq)
	if err != nil {
		return nil, err
	}

	// Copy the update request fields into the existing item
	copier.Copy(existingItem, req)

	if _, err := l.svcCtx.BaseCodeRpcClient.UpdateBaseCodeItem(l.ctx, existingItem); err != nil {
		return nil, err
	}


	return resp, nil
}
