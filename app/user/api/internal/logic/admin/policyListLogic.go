package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PolicyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPolicyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PolicyListLogic {
	return &PolicyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PolicyListLogic) PolicyList(req *types.PolicyItem) (resp *types.PolicyListResp, err error) {
	rpcReq := &pb.BasePolicyItem{
		Domain: req.Domain,
		Role:   req.Role,
		Policy: req.Policy,
		Action: req.Action,
	}
	rpcResp, err := l.svcCtx.UserRpcClient.GetPolicyList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	resp = &types.PolicyListResp{}
	for _, item := range rpcResp.List {
		resp.List = append(resp.List, types.PolicyItem{
			Domain: item.Domain,
			Role:   item.Role,
			Policy: item.Policy,
			Action: item.Action,
		})
	}
	return resp, nil
}
