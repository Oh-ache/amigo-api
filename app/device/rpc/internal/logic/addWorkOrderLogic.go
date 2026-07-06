package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddWorkOrderLogic {
	return &AddWorkOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddWorkOrderLogic) AddWorkOrder(in *pb.AddWorkOrderReq) (*pb.WorkOrderItem, error) {
	m := model.WorkOrder{
		DeviceId: in.DeviceId,
		UserId:   in.UserId,
		Title:    in.Title,
		Content:  in.Content,
		Images:   in.Images,
		Category: in.Category,
		Status:   1,
		IsDelete: 2,
	}

	if _, err := l.svcCtx.WorkOrderModel.Insert(l.ctx, &m); err != nil {
		return nil, err
	}

	return &pb.WorkOrderItem{
		WorkOrderId: m.WorkOrderId,
		DeviceId:    m.DeviceId,
		UserId:      m.UserId,
		Title:       m.Title,
		Content:     m.Content,
		Images:      m.Images,
		Category:    m.Category,
		Status:      m.Status,
		CreateTime:  m.CreateTime,
		UpdateTime:  m.UpdateTime,
	}, nil
}
