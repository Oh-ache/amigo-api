package workOrder

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyLogic {
	return &ReplyLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ReplyLogic) Reply(req *types.ReplyWorkOrderReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.DeviceRpcClient.ReplyWorkOrder(l.ctx, &pb.ReplyWorkOrderReq{
		WorkOrderId: req.WorkOrderId,
		AdminId:     req.AdminId,
		Content:     req.Content,
		Images:      req.Images,
	})
	if err != nil {
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
