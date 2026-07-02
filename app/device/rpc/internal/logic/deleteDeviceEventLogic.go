package logic

import (
	"context"

	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDeviceEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDeviceEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDeviceEventLogic {
	return &DeleteDeviceEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDeviceEventLogic) DeleteDeviceEvent(in *pb.DeleteDeviceEventReq) (*pb.DeleteDeviceEventResp, error) {
	if in.DeviceEventId != 0 {
		if err := l.svcCtx.DeviceEventModel.Delete(l.ctx, in.DeviceEventId); err != nil {
			return nil, err
		}
		return &pb.DeleteDeviceEventResp{Success: true}, nil
	}
	if in.DeviceId != 0 {
		if err := l.svcCtx.DeviceEventModel.DeleteByDevice(l.ctx, in.DeviceId); err != nil {
			return nil, err
		}
		return &pb.DeleteDeviceEventResp{Success: true}, nil
	}
	return &pb.DeleteDeviceEventResp{Success: false}, nil
}
