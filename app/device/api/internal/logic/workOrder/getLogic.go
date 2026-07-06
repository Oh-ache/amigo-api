package workOrder

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetLogic) Get(req *types.GetWorkOrderReq) (resp *types.GetWorkOrderResp, err error) {
	rpcResp, err := l.svcCtx.DeviceRpcClient.GetWorkOrder(l.ctx, &pb.GetWorkOrderReq{
		WorkOrderId: req.WorkOrderId,
	})
	if err != nil {
		return nil, err
	}

	wo := rpcResp.WorkOrder
	resp = &types.GetWorkOrderResp{
		WorkOrder: types.GetWorkOrderItem{
			WorkOrderId: wo.WorkOrderId,
			DeviceId:    wo.DeviceId,
			DeviceName:  wo.DeviceName,
			UserId:      wo.UserId,
			UserName:    wo.UserName,
			Title:       wo.Title,
			Content:     wo.Content,
			Images:      wo.Images,
			Category:    wo.Category,
			Status:      wo.Status,
			CreateTime:  wo.CreateTime,
			UpdateTime:  wo.UpdateTime,
		},
	}

	for _, r := range rpcResp.ReplyList {
		resp.ReplyList = append(resp.ReplyList, types.ReplyItem{
			WorkOrderReplyId: r.WorkOrderReplyId,
			WorkOrderId:      r.WorkOrderId,
			AdminId:          r.AdminId,
			AdminName:        r.AdminName,
			Content:          r.Content,
			Images:           r.Images,
			CreateTime:       r.CreateTime,
		})
	}

	return resp, nil
}
