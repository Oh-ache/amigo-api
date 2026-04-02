package user

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserUpdateReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}

	param := &pb.UpdateUserReq{
		UserId:   req.UserId,
		Username: req.Username,
		Avatar:   req.Avatar,
	}

	_, err = l.svcCtx.UserRpcClient.UpdateUser(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
