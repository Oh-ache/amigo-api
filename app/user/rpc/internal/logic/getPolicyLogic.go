package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPolicyLogic {
	return &GetPolicyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPolicyLogic) GetPolicy(in *pb.BasePolicyItem) (*pb.SuccessResp, error) {
	res := &pb.SuccessResp{}
	res.Success = l.svcCtx.AdminAuth.Enforcer.HasPolicy([]string{in.Role, in.Domain, in.Policy, in.Action})

	return res, nil
}
