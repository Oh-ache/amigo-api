package user

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

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

	// 先尝试通过手机号获取用户
	if req.Mobile != "" {
		getUserReq := &pb.GetUserReq{
			Mobile: req.Mobile,
		}
		userResp, err := l.svcCtx.UserRpcClient.GetUser(l.ctx, getUserReq)
		if err == nil && userResp != nil && userResp.Password == req.Password {
			// 密码匹配，生成登录token
			// 这里可以调用专门的登录RPC来生成token
			// 暂时直接返回用户信息
			resp.UserId = userResp.UserId
			resp.Mobile = userResp.Mobile
			resp.Username = userResp.Username
			resp.Avatar = userResp.Avatar
			resp.CreateTime = int64(userResp.CreateTime)
			resp.Token = "token_placeholder" // 实际应该调用RPC生成token
			return resp, nil
		}
	}

	// 尝试通过用户名获取用户
	if req.Username != "" {
		getUserReq := &pb.GetUserReq{
			Username: req.Username,
		}
		userResp, err := l.svcCtx.UserRpcClient.GetUser(l.ctx, getUserReq)
		if err == nil && userResp != nil && userResp.Password == req.Password {
			resp.UserId = userResp.UserId
			resp.Mobile = userResp.Mobile
			resp.Username = userResp.Username
			resp.Avatar = userResp.Avatar
			resp.CreateTime = int64(userResp.CreateTime)
			resp.Token = "token_placeholder"
			return resp, nil
		}
	}

	return resp, nil
}
