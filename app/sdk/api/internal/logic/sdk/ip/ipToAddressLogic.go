package ip

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IpToAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIpToAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IpToAddressLogic {
	return &IpToAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IpToAddressLogic) IpToAddress(req *types.IpToAddressReq) (resp *types.IpToAddressResp, err error) {
	// todo: add your logic here and delete this line

	return
}
