package app

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
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

func (l *GetLogic) Get(req *types.GetAppReq) (resp *types.GetAppResp, err error) {
	param := &pb.GetAppReq{}
	if err := copier.Copy(param, req); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.DeviceRpcClient.GetApp(l.ctx, param)
	if err != nil {
		return nil, err
	}
	resp = &types.GetAppResp{}
	if err := copier.Copy(resp, rpcResp); err != nil {
		return nil, err
	}
	return resp, nil
}
