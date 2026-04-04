package baseCodeSort

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

)
type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteBaseCodeSortReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	param := &pb.DeleteBaseCodeSortReq{}

	copier.Copy(param, req)
	if _, err := l.svcCtx.BaseCodeRpcClient.DeleteBaseCodeSort(l.ctx, param); err != nil {
		return nil, err
	}


	return resp, nil
}
