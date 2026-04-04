package baseCodeSort

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

)
type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListBaseCodeSortReq) (resp *types.ListBaseCodeSortResp, err error) {
	resp = &types.ListBaseCodeSortResp{}
	param := &pb.ListBaseCodeSortReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.BaseCodeRpcClient.ListBaseCodeSort(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	
	
	return resp, nil
}
