// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.ListUserReq) (resp *types.ListUserResp, err error) {
	rpcReq := &pb.ListUserReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		UserId:   req.UserId,
		Username: req.Username,
		Mobile:   req.Mobile,
	}

	rpcResp, err := l.svcCtx.UserRpcClient.ListUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.ListUserResp{Total: rpcResp.Total}
	for _, item := range rpcResp.List {
		resp.List = append(resp.List, types.GetUserResp{
			UserId:     item.UserId,
			Mobile:     item.Mobile,
			Username:   item.Username,
			Avatar:     item.Avatar,
			CreateTime: int64(item.CreateTime),
		})
	}
	return resp, nil
}
