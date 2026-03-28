package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddPolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPolicyLogic {
	return &AddPolicyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddPolicyLogic) AddPolicy(in *pb.BasePolicyItem) (*pb.SuccessResp, error) {
	res := &pb.SuccessResp{}

	res.Success, _ = l.svcCtx.AdminAuth.Enforcer.AddPolicy([]string{in.Role, in.Domain, in.Policy, in.Action})

	return res, nil
}
