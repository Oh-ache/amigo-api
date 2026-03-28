package logic

import (
	"context"

	"amigo-api/app/user/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MiniappCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMiniappCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MiniappCodeLogic {
	return &MiniappCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MiniappCodeLogic) MiniappCode(in *pb.MiniappCodeReq) (*pb.MiniappCodeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.MiniappCodeResp{}, nil
}
