package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSortLogic {
	return &AddSortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSortLogic) AddSort(req *types.AddBaseCodeSortReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	param := &pb.AddBaseCodeSortReq{}

	copier.Copy(param, req)
	if _, err := l.svcCtx.BaseCodeRpcClient.AddBaseCodeSort(l.ctx, param); err != nil {
		return nil, err
	}

	return resp, nil
}
