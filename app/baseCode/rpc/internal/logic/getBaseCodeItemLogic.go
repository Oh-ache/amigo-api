package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeItemLogic {
	return &GetBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeItemLogic) GetBaseCodeItem(in *pb.GetBaseCodeItemReq) (*pb.BaseCodeItemResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeItemResp{}, nil
}
