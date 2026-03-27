package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeLogic {
	return &GetBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeLogic) GetBaseCode(in *pb.GetBaseCodeReq) (*pb.BaseCodeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeResp{}, nil
}
