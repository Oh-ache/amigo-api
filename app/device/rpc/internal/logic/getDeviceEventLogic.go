package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceEventLogic {
	return &GetDeviceEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceEventLogic) GetDeviceEvent(in *pb.GetDeviceEventReq) (*pb.DeviceEventResp, error) {
	if in.DeviceEventId == 0 {
		return nil, model.ErrNotFound
	}
	m, err := l.svcCtx.DeviceEventModel.FindOne(l.ctx, in.DeviceEventId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return toPbDeviceEvent(m), nil
}
