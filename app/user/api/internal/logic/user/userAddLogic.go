package user

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAddLogic {
	return &UserAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAddLogic) UserAdd(req *types.UserAddReq) (resp *types.GetUserResp, err error) {
	resp = &types.GetUserResp{}

	param := &pb.AddUserReq{
		Username: req.Username,
		Mobile:   req.Mobile,
		Avatar:   req.Avatar,
		Password: req.Password,
	}

	rpcResp, err := l.svcCtx.UserRpcClient.AddUser(l.ctx, param)
	if err != nil {
		return nil, err
	}

	resp.UserId = rpcResp.UserId
	resp.Mobile = rpcResp.Mobile
	resp.Username = rpcResp.Username
	resp.Avatar = rpcResp.Avatar
	return resp, nil
}
