package user

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateMobileLogic {
	return &UserUpdateMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateMobileLogic) UserUpdateMobile(req *types.UserUpdateMobileReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}

	param := &pb.UpdateUserReq{
		UserId: req.UserId,
		Mobile: req.Mobile,
	}

	_, err = l.svcCtx.UserRpcClient.UpdateUser(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
