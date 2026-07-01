package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PolicyAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPolicyAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PolicyAddLogic {
	return &PolicyAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PolicyAddLogic) PolicyAdd(req *types.PolicyItem) (resp *types.EmptyResp, err error) {
	rpcReq := &pb.BasePolicyItem{
		Domain: req.Domain,
		Role:   req.Role,
		Policy: req.Policy,
		Action: req.Action,
	}
	if _, err = l.svcCtx.UserRpcClient.AddPolicy(l.ctx, rpcReq); err != nil {
		l.Errorf("AddPolicy failed: %v", err)
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
