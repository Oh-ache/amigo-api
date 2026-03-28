package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllDomainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllDomainLogic {
	return &GetAllDomainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllDomainLogic) GetAllDomain(in *pb.GetAllDomainReq) (*pb.GetAllDomainResp, error) {
	res := &pb.GetAllDomainResp{}
	res.List, _ = l.svcCtx.AdminAuth.Enforcer.GetAllDomains()

	return res, nil
}
