package workOrder

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
	return &ListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ListLogic) List(req *types.ListWorkOrderReq) (resp *types.ListWorkOrderResp, err error) {
	rpcResp, err := l.svcCtx.DeviceRpcClient.ListWorkOrder(l.ctx, &pb.ListWorkOrderReq{
		DeviceId: req.DeviceId,
		UserId:   req.UserId,
		Status:   req.Status,
		Category: req.Category,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.ListWorkOrderResp{Total: rpcResp.Total}
	for _, item := range rpcResp.List {
		resp.List = append(resp.List, types.GetWorkOrderItem{
			WorkOrderId: item.WorkOrderId,
			DeviceId:    item.DeviceId,
			DeviceName:  item.DeviceName,
			UserId:      item.UserId,
			UserName:    item.UserName,
			Title:       item.Title,
			Content:     item.Content,
			Images:      item.Images,
			Category:    item.Category,
			Status:      item.Status,
			CreateTime:  item.CreateTime,
			UpdateTime:  item.UpdateTime,
		})
	}

	return resp, nil
}
