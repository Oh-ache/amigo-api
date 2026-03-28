package device

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
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

func (l *DeleteLogic) Delete(req *types.DeleteDeviceReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	param := &pb.DeleteDeviceReq{}

	copier.Copy(param, req)
	if _, err := l.svcCtx.DeviceRpcClient.DeleteDevice(l.ctx, param); err != nil {
		return nil, err
	}

	return resp, nil
}
