package logic

import (
	"context"
	"fmt"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCodeLogic {
	return &SendCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendCodeLogic) SendCode(in *pb.SendCodeReq) (*pb.SendCodeResp, error) {
	// TODO 可能要添加发送限制
	resp := &pb.SendCodeResp{}
	ctx := &message.PushContext{}

	if in.Platform == "" {
		in.Platform = "ali_sms"
	}

	ctx.Platform = in.Platform
	ctx.Mobile = in.Mobile

	switch in.Platform {
	case "ali_sms":
		ctx.AliAccessKeySecret = GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.accessKeySecret")
		ctx.AliAccessKeyId = GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.accessKey")
		ctx.TmplateCode = GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.sms.codeTemplateCode")
		ctx.SignName = GetBaseCode(l.ctx, l.svcCtx.BaseCodeRpc, "sdk", "ali.sms.codeSignName")
	}

	code := utils.GetRandomNum()
	content := fmt.Sprintf("{\"code\": \"%s\"}", code)
	ctx.Content = content

	if err := message.PushMessage(ctx); err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("%s%s:%s", utils.SEND_CODE_KEY, in.SendType, ctx.Mobile)
	l.svcCtx.RedisClient.Setex(redisKey, code, 180)

	resp.Success = true

	return resp, nil
}
