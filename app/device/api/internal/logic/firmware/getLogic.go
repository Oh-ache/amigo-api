package firmware

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

func (l *GetLogic) Get(req *types.GetFirmwareReq) (resp *types.GetFirmwareResp, err error) {
	rpcReq := &pb.GetFirmwareReq{
		FirmwareId: req.FirmwareId,
	}

	rpcResp, err := l.svcCtx.DeviceRpcClient.GetFirmware(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.GetFirmwareResp{
		FirmwareId: rpcResp.FirmwareId,
		Name:       rpcResp.Name,
		Version:    rpcResp.Version,
		DeviceType: rpcResp.DeviceType,
		FileUrl:    rpcResp.FileUrl,
		FileSize:   rpcResp.FileSize,
		Md5:        rpcResp.Md5,
		Changelog:  rpcResp.Changelog,
		IsForce:    rpcResp.IsForce,
		IsDelete:   rpcResp.IsDelete,
		CreateTime: rpcResp.CreateTime,
		UpdateTime: rpcResp.UpdateTime,
	}, nil
}
