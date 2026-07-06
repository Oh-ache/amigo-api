package workOrder

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatusLogic {
	return &UpdateStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateStatusLogic) UpdateStatus(req *types.UpdateWorkOrderStatusReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.DeviceRpcClient.UpdateWorkOrderStatus(l.ctx, &pb.UpdateWorkOrderStatusReq{
		WorkOrderId: req.WorkOrderId,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
