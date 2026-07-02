// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package device

import (
	"context"

	"amigo-api/app/device/api/internal/svc"
	"amigo-api/app/device/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AppListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppListLogic {
	return &AppListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppListLogic) AppList(req *types.ListAppReq) (resp *types.ListAppResp, err error) {
	resp = &types.ListAppResp{}
	param := &pb.ListAppReq{}
	if err := copier.Copy(param, req); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.DeviceRpcClient.ListApp(l.ctx, param)
	if err != nil {
		return nil, err
	}
	if err := copier.Copy(resp, rpcResp); err != nil {
		return nil, err
	}
	return resp, nil
}
