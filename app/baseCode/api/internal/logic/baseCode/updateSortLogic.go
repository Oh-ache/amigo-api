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
	resp = &types.EmptyResp{}
	// First, get the existing sort to populate the BaseCodeSortResp
	getReq := &pb.GetBaseCodeSortReq{
		BaseCodeSortId: req.BaseCodeSortId,
	}
	existingSort, err := l.svcCtx.BaseCodeRpcClient.GetBaseCodeSort(l.ctx, getReq)
	if err != nil {
		return nil, err
	}

	// Copy the update request fields into the existing item
	copier.Copy(existingSort, req)

	if _, err := l.svcCtx.BaseCodeRpcClient.UpdateBaseCodeSort(l.ctx, existingSort); err != nil {
		return nil, err
	}

	return resp, nil
}
