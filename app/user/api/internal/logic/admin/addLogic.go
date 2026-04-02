package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.AdminAddReq) (resp *types.GetAdminResp, err error) {
	resp = &types.GetAdminResp{}

	param := &pb.AddAdminReq{
		Mobile:   req.Mobile,
		Username: req.Username,
		Avatar:   req.Avatar,
		Password: req.Password,
	}

	_, err = l.svcCtx.UserRpcClient.AddAdmin(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
