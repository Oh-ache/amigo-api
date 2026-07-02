// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListAdminReq) (resp *types.ListAdminResp, err error) {
	rpcReq := &pb.ListAdminReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		AdminId:  req.AdminId,
		Username: req.Username,
		Mobile:   req.Mobile,
	}

	rpcResp, err := l.svcCtx.UserRpcClient.ListAdmin(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.ListAdminResp{Total: rpcResp.Total}
	for _, item := range rpcResp.List {
		resp.List = append(resp.List, types.GetAdminResp{
			AdminId:    item.AdminId,
			Mobile:     item.Mobile,
			Username:   item.Username,
			Avatar:     item.Avatar,
			CreateTime: int64(item.CreateTime),
		})
	}
	return resp, nil
}
