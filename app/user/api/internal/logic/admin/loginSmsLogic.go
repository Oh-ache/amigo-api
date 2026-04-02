package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginSmsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginSmsLogic {
	return &LoginSmsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginSmsLogic) LoginSms(req *types.AdminLoginSmsReq) (resp *types.AdminLoginResp, err error) {
	resp = &types.AdminLoginResp{}

	// 短信登录逻辑 - 这里可以调用相关的RPC
	// 暂时留空，因为proto中可能没有专门的短信登录方法
	// 可以复用LoginAdmin或者添加新的RPC方法

	return resp, nil
}
