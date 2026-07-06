package workOrder

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddLogic) Add(req *types.AddWorkOrderReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.DeviceRpcClient.AddWorkOrder(l.ctx, &pb.AddWorkOrderReq{
		DeviceId: req.DeviceId,
		UserId:   req.UserId,
		Title:    req.Title,
		Content:  req.Content,
		Images:   req.Images,
		Category: req.Category,
	})
	if err != nil {
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
