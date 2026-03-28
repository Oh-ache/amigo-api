package device

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

func (l *GetLogic) Get(req *types.GetDeviceReq) (resp *types.GetDeviceResp, err error) {
	resp = &types.GetDeviceResp{}
	param := &pb.GetDeviceReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.DeviceRpcClient.GetDevice(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
