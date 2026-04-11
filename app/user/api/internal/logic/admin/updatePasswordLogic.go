package admin

import (
	"context"
	"fmt"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(req *types.AdminUpdatePasswordReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}

	// 验证新密码和确认密码一致
	if req.Password != req.RePassword {
		return nil, fmt.Errorf("两次密码输入不一致")
	}

	param := &pb.UpdateAdminReq{
		AdminId:  req.AdminId,
		Password: req.Password,
	}

	_, err = l.svcCtx.UserRpcClient.UpdateAdmin(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
