package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePolicyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePolicyLogic {
	return &DeletePolicyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePolicyLogic) DeletePolicy(in *pb.BasePolicyItem) (*pb.SuccessResp, error) {
	res := &pb.SuccessResp{}

	res.Success, _ = l.svcCtx.AdminAuth.Enforcer.RemovePolicy(in.Role, in.Domain, in.Policy, in.Action)

	return res, nil
}
