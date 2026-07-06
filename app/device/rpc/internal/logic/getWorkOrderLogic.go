package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkOrderLogic {
	return &GetWorkOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetWorkOrderLogic) GetWorkOrder(in *pb.GetWorkOrderReq) (*pb.GetWorkOrderResp, error) {
	wo, err := l.svcCtx.WorkOrderModel.FindOne(l.ctx, in.WorkOrderId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	var deviceName string
	if dev, devErr := l.svcCtx.DeviceModel.FindOne(l.ctx, wo.DeviceId); devErr == nil {
		deviceName = dev.Name
	}

	replyList, err := l.svcCtx.WorkOrderReplyModel.ListByWorkOrderId(l.ctx, wo.WorkOrderId)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetWorkOrderResp{
		WorkOrder: &pb.WorkOrderItem{
			WorkOrderId: wo.WorkOrderId,
			DeviceId:    wo.DeviceId,
			DeviceName:  deviceName,
			UserId:      wo.UserId,
			Title:       wo.Title,
			Content:     wo.Content,
			Images:      wo.Images,
			Category:    wo.Category,
			Status:      wo.Status,
			CreateTime:  wo.CreateTime,
			UpdateTime:  wo.UpdateTime,
		},
	}

	for _, r := range replyList {
		resp.ReplyList = append(resp.ReplyList, &pb.ReplyItem{
			WorkOrderReplyId: r.WorkOrderReplyId,
			WorkOrderId:      r.WorkOrderId,
			AdminId:          r.AdminId,
			Content:          r.Content,
			Images:           r.Images,
			CreateTime:       r.CreateTime,
		})
	}

	return resp, nil
}
