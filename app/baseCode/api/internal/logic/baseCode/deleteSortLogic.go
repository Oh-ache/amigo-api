package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSortLogic {
	return &DeleteSortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSortLogic) DeleteSort(req *types.DeleteBaseCodeSortReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	param := &pb.DeleteBaseCodeSortReq{}

	copier.Copy(param, req)
	if _, err := l.svcCtx.BaseCodeRpcClient.DeleteBaseCodeSort(l.ctx, param); err != nil {
		return nil, err
	}

	return resp, nil
}
