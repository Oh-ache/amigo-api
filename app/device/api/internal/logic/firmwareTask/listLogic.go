package firmwareTask

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

func (l *ListLogic) List(req *types.ListFirmwareTaskReq) (resp *types.ListFirmwareTaskResp, err error) {
	rpcReq := &pb.ListFirmwareTaskReq{
		FirmwareId: req.FirmwareId,
		DeviceId:   req.DeviceId,
		Status:     req.Status,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	rpcResp, err := l.svcCtx.DeviceRpcClient.ListFirmwareTask(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.ListFirmwareTaskResp{Total: rpcResp.Total}
	for _, item := range rpcResp.List {
		resp.List = append(resp.List, types.GetFirmwareTaskResp{
			FirmwareTaskId: item.FirmwareTaskId,
			FirmwareId:     item.FirmwareId,
			DeviceId:       item.DeviceId,
			Status:         item.Status,
			Progress:       item.Progress,
			ErrorMsg:       item.ErrorMsg,
			StartedAt:      item.StartedAt,
			CompletedAt:    item.CompletedAt,
			CreateTime:     item.CreateTime,
			UpdateTime:     item.UpdateTime,
		})
	}

	return resp, nil
}
