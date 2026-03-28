package ip

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"
	"amigo-api/app/sdk/rpc/sdk"

	"github.com/jinzhu/copier"
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
	resp = &types.IpToAddressResp{}
	param := &sdk.IpToAddressReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.SdkRpcClient.IpToAddress(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
