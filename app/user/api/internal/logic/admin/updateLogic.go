package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.AdminUpdateReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}

	param := &pb.UpdateAdminReq{
		AdminId:  req.AdminId,
		Username: req.Username,
		Avatar:   req.Avatar,
	}

	_, err = l.svcCtx.UserRpcClient.UpdateAdmin(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
