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

type EventListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EventListLogic {
	return &EventListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventListLogic) EventList(req *types.ListDeviceEventReq) (resp *types.ListDeviceEventResp, err error) {
	resp = &types.ListDeviceEventResp{}
	param := &pb.ListDeviceEventReq{}
	if err := copier.Copy(param, req); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.DeviceRpcClient.ListDeviceEvent(l.ctx, param)
	if err != nil {
		return nil, err
	}
	if err := copier.Copy(resp, rpcResp); err != nil {
		return nil, err
	}
	return resp, nil
}
