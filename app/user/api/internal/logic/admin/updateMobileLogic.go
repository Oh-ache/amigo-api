package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMobileLogic {
	return &UpdateMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMobileLogic) UpdateMobile(req *types.AdminUpdateMobileReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}

	param := &pb.UpdateAdminReq{
		AdminId: req.AdminId,
		Mobile:  req.Mobile,
	}

	_, err = l.svcCtx.UserRpcClient.UpdateAdmin(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
