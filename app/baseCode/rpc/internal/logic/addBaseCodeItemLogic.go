package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBaseCodeItemLogic {
	return &AddBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBaseCodeItemLogic) AddBaseCodeItem(in *pb.AddBaseCodeItemReq) (*pb.BaseCodeItemResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeItemResp{}, nil
}
