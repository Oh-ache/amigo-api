package admin

import (
	"context"

	"amigo-api/app/user/api/internal/svc"
	"amigo-api/app/user/api/internal/types"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DomainListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDomainListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DomainListLogic {
	return &DomainListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DomainListLogic) DomainList() (resp *types.DomainListResp, err error) {
	rpcResp, err := l.svcCtx.UserRpcClient.GetAllDomain(l.ctx, &pb.GetAllDomainReq{})
	if err != nil {
		l.Errorf("GetAllDomain failed: %v", err)
		return nil, err
	}
	resp = &types.DomainListResp{
		List: rpcResp.List,
	}
	return resp, nil
}
