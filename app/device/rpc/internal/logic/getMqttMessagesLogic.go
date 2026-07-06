package logic

import (
	"context"
	"encoding/json"

	"amigo-api/app/device/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMqttMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMqttMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMqttMessagesLogic {
	return &GetMqttMessagesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetMqttMessagesLogic) GetMqttMessages(in *pb.GetMqttMessagesReq) (*pb.GetMqttMessagesResp, error) {
	limit := in.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	vals, err := l.svcCtx.RedisClient.Lrange("mqtt:messages:latest", 0, int(limit-1))
	if err != nil {
		return nil, err
	}

	resp := &pb.GetMqttMessagesResp{}
	for _, v := range vals {
		var msg svc.MQTTMessage
		if err := json.Unmarshal([]byte(v), &msg); err != nil {
			continue
		}
		resp.List = append(resp.List, &pb.MqttMsg{
			Topic:     msg.Topic,
			Payload:   msg.Payload,
			Timestamp: msg.Timestamp,
		})
	}

	return resp, nil
}
