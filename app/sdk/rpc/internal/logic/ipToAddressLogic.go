package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"

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
	// todo: add your logic here and delete this line

	return &pb.IpToAddressResp{}, nil
}
