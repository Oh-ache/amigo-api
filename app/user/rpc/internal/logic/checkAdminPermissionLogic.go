package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckAdminPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckAdminPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckAdminPermissionLogic {
	return &CheckAdminPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckAdminPermissionLogic) CheckAdminPermission(in *pb.CheckAdminPermissionReq) (*pb.SuccessResp, error) {
	res := &pb.SuccessResp{}

	res.Success, _ = l.svcCtx.AdminAuth.Enforcer.Enforce(in.AdminId, in.Domain, in.Policy, in.Action)

	return res, nil
}
