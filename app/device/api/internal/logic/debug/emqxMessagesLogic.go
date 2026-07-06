package debug

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmqxMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmqxMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmqxMessagesLogic {
	return &EmqxMessagesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *EmqxMessagesLogic) EmqxMessages(req *types.GetMqttMessagesReq) (resp *types.GetMqttMessagesResp, err error) {
	rpcResp, err := l.svcCtx.DeviceRpcClient.GetMqttMessages(l.ctx, &pb.GetMqttMessagesReq{
		Limit: req.Limit,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetMqttMessagesResp{}
	for _, item := range rpcResp.List {
		resp.List = append(resp.List, types.MqttMsg{
			Topic:     item.Topic,
			Payload:   item.Payload,
			Timestamp: item.Timestamp,
		})
	}

	return resp, nil
}
