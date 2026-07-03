package firmwareTask

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.GetFirmwareTaskReq) (resp *types.GetFirmwareTaskResp, err error) {
	rpcReq := &pb.GetFirmwareTaskReq{
		FirmwareTaskId: req.FirmwareTaskId,
	}

	rpcResp, err := l.svcCtx.DeviceRpcClient.GetFirmwareTask(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.GetFirmwareTaskResp{
		FirmwareTaskId: rpcResp.FirmwareTaskId,
		FirmwareId:     rpcResp.FirmwareId,
		DeviceId:       rpcResp.DeviceId,
		Status:         rpcResp.Status,
		Progress:       rpcResp.Progress,
		ErrorMsg:       rpcResp.ErrorMsg,
		StartedAt:      rpcResp.StartedAt,
		CompletedAt:    rpcResp.CompletedAt,
		CreateTime:     rpcResp.CreateTime,
		UpdateTime:     rpcResp.UpdateTime,
	}, nil
}
