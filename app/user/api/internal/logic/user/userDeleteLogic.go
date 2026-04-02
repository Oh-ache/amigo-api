package user

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteLogic) UserDelete(req *types.UserDeleteReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}

	param := &pb.DeleteUserReq{
		UserId: req.UserId,
	}

	_, err = l.svcCtx.UserRpcClient.DeleteUser(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
