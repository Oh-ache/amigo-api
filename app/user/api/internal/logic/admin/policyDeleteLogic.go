package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PolicyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPolicyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PolicyDeleteLogic {
	return &PolicyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PolicyDeleteLogic) PolicyDelete(req *types.PolicyItem) (resp *types.EmptyResp, err error) {
	rpcReq := &pb.BasePolicyItem{
		Domain: req.Domain,
		Role:   req.Role,
		Policy: req.Policy,
		Action: req.Action,
	}
	if _, err = l.svcCtx.UserRpcClient.DeletePolicy(l.ctx, rpcReq); err != nil {
		l.Errorf("DeletePolicy failed: %v", err)
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
