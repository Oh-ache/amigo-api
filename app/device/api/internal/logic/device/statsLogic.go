package device

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsLogic {
	return &StatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsLogic) Stats() (resp *types.DeviceStatsResp, err error) {
	rpcResp, err := l.svcCtx.DeviceRpcClient.GetDeviceStats(l.ctx, &pb.GetDeviceStatsReq{})
	if err != nil {
		return nil, err
	}

	return &types.DeviceStatsResp{
		Total:   rpcResp.Total,
		Online:  rpcResp.Online,
		Offline: rpcResp.Offline,
		Warning: rpcResp.Warning,
	}, nil
}
