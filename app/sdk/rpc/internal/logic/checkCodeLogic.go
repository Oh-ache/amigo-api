package logic

import (
	"context"
	"fmt"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckCodeLogic {
	return &CheckCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckCodeLogic) CheckCode(in *pb.CheckCodeReq) (*pb.CheckCodeResp, error) {
	if in.Platform == "" {
		in.Platform = "ali_sms"
	}
	resp := &pb.CheckCodeResp{}
	redisKey := fmt.Sprintf("%s%s:%s", utils.SEND_CODE_KEY, in.SendType, in.Mobile)

	if value, err := l.svcCtx.RedisClient.Get(redisKey); err != nil || len(value) == 0 {
		return nil, fmt.Errorf("请发送验证码")
	} else {
		if value != in.Code {
			return nil, fmt.Errorf("验证码错误")
		}
	}

	resp.Success = true
	_, _ = l.svcCtx.RedisClient.Del(redisKey)
	return resp, nil
}
