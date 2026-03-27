package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBaseCodeLogic {
	return &AddBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBaseCodeLogic) AddBaseCode(in *pb.AddBaseCodeReq) (*pb.BaseCodeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeResp{}, nil
}
