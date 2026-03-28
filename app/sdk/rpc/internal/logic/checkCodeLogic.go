package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckCodeLogic {
	return &CheckCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckCodeLogic) CheckCode(in *pb.CheckCodeReq) (*pb.CheckCodeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CheckCodeResp{}, nil
}
