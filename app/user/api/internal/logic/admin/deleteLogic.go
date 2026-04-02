package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.AdminDeleteReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}

	param := &pb.DeleteAdminReq{
		AdminId: req.AdminId,
	}

	_, err = l.svcCtx.UserRpcClient.DeleteAdmin(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
