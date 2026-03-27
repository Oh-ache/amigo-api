package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeLogic {
	return &ListBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeLogic) ListBaseCode(in *pb.ListBaseCodeReq) (*pb.ListBaseCodeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListBaseCodeResp{}, nil
}
