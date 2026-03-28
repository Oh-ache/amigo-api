package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleListLogic {
	return &GetRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleListLogic) GetRoleList(in *pb.BaseRoleItem) (*pb.GetRoleListeResp, error) {
	res := &pb.GetRoleListeResp{}

	if len(in.Domain) > 0 && len(in.AdminId) == 0 {
		list := l.svcCtx.AdminAuth.Enforcer.GetAllRolesByDomain(in.Domain)
		list = utils.RemoveItem(list, "_")
		for _, item := range list {
			baseRole := &pb.BaseRoleItem{
				Domain: in.Domain,
				Role:   item,
			}

			res.List = append(res.List, baseRole)
		}
	} else {
		list, _ := l.svcCtx.AdminAuth.Enforcer.GetRolesForUser(in.AdminId, in.Domain)
		list = utils.RemoveItem(list, "_")
		for _, item := range list {
			baseRole := &pb.BaseRoleItem{
				Domain:  in.Domain,
				Role:    item,
				AdminId: in.AdminId,
			}

			res.List = append(res.List, baseRole)
		}
	}

	return res, nil
}
