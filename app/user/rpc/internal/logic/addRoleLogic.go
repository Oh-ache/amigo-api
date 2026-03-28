package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddRoleLogic) AddRole(in *pb.BaseRoleItem) (*pb.SuccessResp, error) {
	res := &pb.SuccessResp{}

	if len(in.Domain) > 0 {
		res.Success, _ = l.svcCtx.AdminAuth.Enforcer.AddRoleForUserInDomain(in.AdminId, in.Role, in.Domain)
	}

	return res, nil
}
