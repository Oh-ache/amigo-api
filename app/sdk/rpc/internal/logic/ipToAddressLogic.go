package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"
	"amigo-api/common/utils/plug/ip"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type IpToAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIpToAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IpToAddressLogic {
	return &IpToAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IpToAddressLogic) IpToAddress(in *pb.IpToAddressReq) (*pb.IpToAddressResp, error) {
	resp := &pb.IpToAddressResp{}
	param := &ip.Ip2AddressReq{}
	param.Ip = in.Ip

	getBaseCodeReq := &pb.GetBaseCodeReq{}
	getBaseCodeReq.SortKey = "sdk"
	item := &pb.BaseCodeResp{}

	getBaseCodeReq.Key = "ip.tianyan.appid"
	item, _ = l.svcCtx.BaseCodeRpc.GetBaseCode(l.ctx, getBaseCodeReq)
	param.AppId = item.Content

	getBaseCodeReq.Key = "ip.tianyan.appsecurity"
	item, _ = l.svcCtx.BaseCodeRpc.GetBaseCode(l.ctx, getBaseCodeReq)
	param.AppSecurity = item.Content

	if result, err := ip.Ip2Address(param); err != nil {
		return nil, err
	} else {
		copier.Copy(resp, result)
	}

	return resp, nil
}
