package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAddLogic {
	return &RoleAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAddLogic) RoleAdd(req *types.RoleItem) (resp *types.EmptyResp, err error) {
	rpcReq := &pb.BaseRoleItem{
		Domain:  req.Domain,
		Role:    req.Role,
		AdminId: req.AdminId,
	}
	if _, err = l.svcCtx.UserRpcClient.AddRole(l.ctx, rpcReq); err != nil {
		l.Errorf("AddRole failed: %v", err)
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
