package firmware

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

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

func (l *UpdateLogic) Update(req *types.UpdateFirmwareReq) (resp *types.EmptyResp, err error) {
	rpcReq := &pb.FirmwareResp{
		FirmwareId: req.FirmwareId,
		Name:       req.Name,
		Changelog:  req.Changelog,
		IsForce:    req.IsForce,
	}

	if _, err := l.svcCtx.DeviceRpcClient.UpdateFirmware(l.ctx, rpcReq); err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
