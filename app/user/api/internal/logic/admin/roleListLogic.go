package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleListLogic {
	return &RoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleListLogic) RoleList(req *types.RoleItem) (resp *types.RoleListResp, err error) {
	rpcReq := &pb.BaseRoleItem{
		Domain:  req.Domain,
		Role:    req.Role,
		AdminId: req.AdminId,
	}
	rpcResp, err := l.svcCtx.UserRpcClient.GetRoleList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp = &types.RoleListResp{}
	for _, item := range rpcResp.List {
		resp.List = append(resp.List, types.RoleItem{
			Domain:  item.Domain,
			Role:    item.Role,
			AdminId: item.AdminId,
		})
	}
	return resp, nil
}
