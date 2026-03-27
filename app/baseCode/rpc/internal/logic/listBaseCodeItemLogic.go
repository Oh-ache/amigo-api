package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeItemLogic {
	return &ListBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeItemLogic) ListBaseCodeItem(in *pb.ListBaseCodeItemReq) (*pb.ListBaseCodeItemResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListBaseCodeItemResp{}, nil
}
