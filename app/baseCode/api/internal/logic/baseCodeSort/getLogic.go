package baseCodeSort

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

)
type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.GetBaseCodeSortReq) (resp *types.GetBaseCodeSortResp, err error) {
	resp = &types.GetBaseCodeSortResp{}
	param := &pb.GetBaseCodeSortReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.BaseCodeRpcClient.GetBaseCodeSort(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	
	
	return resp, nil
}
