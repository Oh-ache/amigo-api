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

	if len(in.Domain) > 0 && len(in.Role) > 0 {
		list := l.svcCtx.AdminAuth.Enforcer.GetFilteredPolicy(1, in.Domain)
		for _, item := range list {
			if item[0] == in.Role {
				res.List = append(res.List, &pb.BasePolicyItem{
					Role:   item[0],
					Domain: item[1],
					Policy: item[2],
					Action: item[3],
				})
			}
		}
		return res, nil
	}

	var list [][]string
	if len(in.Domain) > 0 {
		list = l.svcCtx.AdminAuth.Enforcer.GetFilteredPolicy(1, in.Domain)
	} else if len(in.Role) > 0 {
		list = l.svcCtx.AdminAuth.Enforcer.GetFilteredPolicy(0, in.Role)
	} else {
		list = l.svcCtx.AdminAuth.Enforcer.GetPolicy()
	}

	for _, item := range list {
		res.List = append(res.List, &pb.BasePolicyItem{
			Role:   item[0],
			Domain: item[1],
			Policy: item[2],
			Action: item[3],
		})
	}

	return res, nil
}
