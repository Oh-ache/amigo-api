package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseCodeItemLogic {
	return &UpdateBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBaseCodeItemLogic) UpdateBaseCodeItem(in *pb.BaseCodeItemResp) (*pb.BaseCodeItemResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeItemResp{}, nil
}
