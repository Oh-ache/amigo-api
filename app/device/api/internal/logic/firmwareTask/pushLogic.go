package firmwareTask

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPushLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushLogic {
	return &PushLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PushLogic) Push(req *types.PushFirmwareTaskReq) (resp *types.EmptyResp, err error) {
	rpcReq := &pb.PushFirmwareTaskReq{
		FirmwareId: req.FirmwareId,
		DeviceIds:  req.DeviceIds,
	}

	if _, err := l.svcCtx.DeviceRpcClient.PushFirmwareTask(l.ctx, rpcReq); err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
