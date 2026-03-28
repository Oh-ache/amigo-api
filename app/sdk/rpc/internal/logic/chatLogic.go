package logic

import (
	"context"
	"fmt"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatLogic) Chat(in *pb.ChatReq) (*pb.ChatResp, error) {
	resp := &pb.ChatResp{}
	redisKey := fmt.Sprintf("%s%s", utils.CHAT_KEY, in.Content)
	if value, _ := l.svcCtx.RedisClient.Get(redisKey); value != "" {
		resp.Content = value
		return resp, nil
	}

	model := GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "xinghuo.chat.model")
	passwd := GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "xinghuo.chat.apipasswd")

	answer, err := chat.XinghuoChat(passwd, model, in.Content)

	l.svcCtx.RedisClient.Set(redisKey, answer)
	l.svcCtx.RedisClient.Expire(redisKey, 180)

	if err != nil {
		resp.Content = ""
		return resp, err
	}
	resp.Content = answer
	return resp, nil
}
