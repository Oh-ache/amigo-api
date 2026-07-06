package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkOrderLogic {
	return &ListWorkOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListWorkOrderLogic) ListWorkOrder(in *pb.ListWorkOrderReq) (*pb.ListWorkOrderResp, error) {
	search := &model.WorkOrderSearch{
		DeviceId: in.DeviceId,
		UserId:   in.UserId,
		Status:   in.Status,
		Category: in.Category,
		Page:     in.Page,
		PageSize: in.PageSize,
	}

	list, total, err := l.svcCtx.WorkOrderModel.List(l.ctx, search)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListWorkOrderResp{Total: total}
	for _, item := range list {
		var deviceName string
		if dev, devErr := l.svcCtx.DeviceModel.FindOne(l.ctx, item.DeviceId); devErr == nil {
			deviceName = dev.Name
		}

		resp.List = append(resp.List, &pb.WorkOrderItem{
			WorkOrderId: item.WorkOrderId,
			DeviceId:    item.DeviceId,
			DeviceName:  deviceName,
			UserId:      item.UserId,
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
