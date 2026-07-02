// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package device

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type EventGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EventGetLogic {
	return &EventGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventGetLogic) EventGet(req *types.GetDeviceEventReq) (resp *types.GetDeviceEventResp, err error) {
	param := &pb.GetDeviceEventReq{}
	if err := copier.Copy(param, req); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.DeviceRpcClient.GetDeviceEvent(l.ctx, param)
	if err != nil {
		return nil, err
	}
	resp = &types.GetDeviceEventResp{}
	if err := copier.Copy(resp, rpcResp); err != nil {
		return nil, err
	}
	return resp, nil
}
