package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPolicyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPolicyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPolicyListLogic {
	return &GetPolicyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPolicyListLogic) GetPolicyList(in *pb.BasePolicyItem) (*pb.GetPolicyListResp, error) {
	res := &pb.GetPolicyListResp{}
	list := make([][]string, 0)

	if len(in.Domain) == 0 {
		list = l.svcCtx.AdminAuth.Enforcer.GetFilteredPolicy(1, in.Domain)
	}

	if len(in.Role) > 0 {
		list = l.svcCtx.AdminAuth.Enforcer.GetFilteredPolicy(0, in.Role)
	}

	for _, item := range list {
		policy := &pb.BasePolicyItem{
			Role:   item[0],
			Domain: item[1],
			Policy: item[2],
			Action: item[3],
		}
		res.List = append(res.List, policy)
	}

	return res, nil
}
