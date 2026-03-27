package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBaseCodeLogic {
	return &DeleteBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBaseCodeLogic) DeleteBaseCode(in *pb.DeleteBaseCodeReq) (*pb.DeleteBaseCodeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DeleteBaseCodeResp{}, nil
}
