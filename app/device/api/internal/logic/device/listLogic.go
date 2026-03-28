package device

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

func (l *ListLogic) List(req *types.ListDeviceReq) (resp *types.ListDeviceResp, err error) {
	resp = &types.ListDeviceResp{}
	param := &pb.ListDeviceReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.DeviceRpcClient.ListDevice(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
