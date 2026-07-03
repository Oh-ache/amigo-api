package firmware

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListFirmwareReq) (resp *types.ListFirmwareResp, err error) {
	rpcReq := &pb.ListFirmwareReq{
		DeviceType: req.DeviceType,
		Version:    req.Version,
		Name:       req.Name,
		IsForce:    req.IsForce,
		IsDelete:   req.IsDelete,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	rpcResp, err := l.svcCtx.DeviceRpcClient.ListFirmware(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.ListFirmwareResp{Total: rpcResp.Total}
	for _, item := range rpcResp.List {
		resp.List = append(resp.List, types.GetFirmwareResp{
			FirmwareId: item.FirmwareId,
			Name:       item.Name,
			Version:    item.Version,
			DeviceType: item.DeviceType,
			FileUrl:    item.FileUrl,
			FileSize:   item.FileSize,
			Md5:        item.Md5,
			Changelog:  item.Changelog,
			IsForce:    item.IsForce,
			IsDelete:   item.IsDelete,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}

	return resp, nil
}
