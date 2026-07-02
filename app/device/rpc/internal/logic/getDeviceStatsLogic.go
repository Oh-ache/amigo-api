package logic

import (
	"context"

	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceStatsLogic {
	return &GetDeviceStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceStatsLogic) GetDeviceStats(in *pb.GetDeviceStatsReq) (*pb.DeviceStatsResp, error) {
	stats, err := l.svcCtx.DeviceModel.Stats(l.ctx)
	if err != nil {
		return nil, err
	}

	return &pb.DeviceStatsResp{
		Total:   stats.Total,
		Online:  stats.Online,
		Offline: stats.Offline,
		Warning: stats.Warning,
	}, nil
}
