package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseCodeLogic {
	return &UpdateBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBaseCodeLogic) UpdateBaseCode(in *pb.BaseCodeResp) (*pb.BaseCodeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeResp{}, nil
}
