package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDeleteLogic {
	return &RoleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleDeleteLogic) RoleDelete(req *types.RoleItem) (resp *types.EmptyResp, err error) {
	rpcReq := &pb.BaseRoleItem{
		Domain:  req.Domain,
		Role:    req.Role,
		AdminId: req.AdminId,
	}
	if _, err = l.svcCtx.UserRpcClient.DeleteRole(l.ctx, rpcReq); err != nil {
		l.Errorf("DeleteRole failed: %v", err)
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
