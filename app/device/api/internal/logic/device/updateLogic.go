package device

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateDeviceReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	getReq := &pb.GetDeviceReq{
		DeviceId: req.DeviceId,
	}
	existingItem, err := l.svcCtx.DeviceRpcClient.GetDevice(l.ctx, getReq)
	if err != nil {
		return nil, err
	}

	copier.Copy(existingItem, req)

	if _, err := l.svcCtx.DeviceRpcClient.UpdateDevice(l.ctx, existingItem); err != nil {
		return nil, err
	}

	return resp, nil
}
