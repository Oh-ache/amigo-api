package user

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.UserLoginReq) (resp *types.UserLoginResp, err error) {
	resp = &types.UserLoginResp{}

	// 使用 username 登录
	param := &pb.LoginResp{
		Username: req.Username,
		Password: req.Password,
	}

	rpcResp, err := l.svcCtx.UserRpcClient.Login(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
