package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyWorkOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyWorkOrderLogic {
	return &ReplyWorkOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyWorkOrderLogic) ReplyWorkOrder(in *pb.ReplyWorkOrderReq) (*pb.ReplyItem, error) {
	_, err := l.svcCtx.WorkOrderModel.FindOne(l.ctx, in.WorkOrderId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	reply := model.WorkOrderReply{
		WorkOrderId: in.WorkOrderId,
		AdminId:     in.AdminId,
		Content:     in.Content,
		Images:      in.Images,
		IsDelete:    2,
	}

	if _, err := l.svcCtx.WorkOrderReplyModel.Insert(l.ctx, &reply); err != nil {
		return nil, err
	}

	// 更新工单状态为"已回复"
	if current, err := l.svcCtx.WorkOrderModel.FindOne(l.ctx, in.WorkOrderId); err == nil {
		current.Status = 3
		l.svcCtx.WorkOrderModel.Update(l.ctx, current)
	}

	return &pb.ReplyItem{
		WorkOrderReplyId: reply.WorkOrderReplyId,
		WorkOrderId:      reply.WorkOrderId,
		AdminId:          reply.AdminId,
		Content:          reply.Content,
		Images:           reply.Images,
		CreateTime:       reply.CreateTime,
	}, nil
}
