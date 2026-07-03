package firmware

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.AddFirmwareReq) (resp *types.EmptyResp, err error) {
	rpcReq := &pb.AddFirmwareReq{
		Name:       req.Name,
		Version:    req.Version,
		DeviceType: req.DeviceType,
		FileUrl:    req.FileUrl,
		FileSize:   req.FileSize,
		Md5:        req.Md5,
		Changelog:  req.Changelog,
		IsForce:    req.IsForce,
	}

	if _, err := l.svcCtx.DeviceRpcClient.AddFirmware(l.ctx, rpcReq); err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
