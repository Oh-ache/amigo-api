package firmware

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

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

func (l *DeleteLogic) Delete(req *types.DeleteFirmwareReq) (resp *types.EmptyResp, err error) {
	rpcReq := &pb.DeleteFirmwareReq{
		FirmwareId: req.FirmwareId,
	}

	if _, err := l.svcCtx.DeviceRpcClient.DeleteFirmware(l.ctx, rpcReq); err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
