package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoadUserPolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoadUserPolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoadUserPolicyLogic {
	return &LoadUserPolicyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoadUserPolicyLogic) LoadUserPolicy(in *pb.BaseRoleItem) (*pb.SuccessResp, error) {
	filter := xormadapter.Filter{
		Ptype: []string{"g"},
		V0:    []string{in.AdminId},
		V1:    []string{in.Domain},
	}

	res := &pb.SuccessResp{}

	if err := l.svcCtx.AdminAuth.Enforcer.LoadFilteredPolicy(filter); err != nil {
		return res, err
	}

	res.Success = true
	return res, nil
}
