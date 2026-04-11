package user

import (
	"context"
	"fmt"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdatePasswordLogic {
	return &UserUpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdatePasswordLogic) UserUpdatePassword(req *types.UserUpdatePasswordReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}

	// 验证新密码和确认密码一致
	if req.Password != req.RePassword {
		return nil, fmt.Errorf("两次密码输入不一致")
	}

	param := &pb.UpdateUserReq{
		UserId:   req.UserId,
		Password: req.Password,
	}

	_, err = l.svcCtx.UserRpcClient.UpdateUser(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
