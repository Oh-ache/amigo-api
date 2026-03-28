package device

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.AddDeviceReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	param := &pb.AddDeviceReq{}

	copier.Copy(param, req)
	if _, err := l.svcCtx.DeviceRpcClient.AddDevice(l.ctx, param); err != nil {
		return nil, err
	}

	return resp, nil
}
