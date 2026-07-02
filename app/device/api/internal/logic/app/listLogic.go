package app

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
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

func (l *ListLogic) List(req *types.ListAppReq) (resp *types.ListAppResp, err error) {
	resp = &types.ListAppResp{}
	param := &pb.ListAppReq{}
	if err := copier.Copy(param, req); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.DeviceRpcClient.ListApp(l.ctx, param)
	if err != nil {
		return nil, err
	}
	if err := copier.Copy(resp, rpcResp); err != nil {
		return nil, err
	}
	return resp, nil
}
