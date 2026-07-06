package logic

import (
	"context"

	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkOrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWorkOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkOrderStatusLogic {
	return &UpdateWorkOrderStatusLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateWorkOrderStatusLogic) UpdateWorkOrderStatus(in *pb.UpdateWorkOrderStatusReq) (*pb.WorkOrderStatusResp, error) {
	wo, err := l.svcCtx.WorkOrderModel.FindOne(l.ctx, in.WorkOrderId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	wo.Status = in.Status

	if err := l.svcCtx.WorkOrderModel.Update(l.ctx, wo); err != nil {
		return nil, err
	}

	return &pb.WorkOrderStatusResp{Success: true}, nil
}
