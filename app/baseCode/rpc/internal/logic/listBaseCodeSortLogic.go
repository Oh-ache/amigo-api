package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBaseCodeSortLogic {
	return &ListBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListBaseCodeSortLogic) ListBaseCodeSort(in *pb.ListBaseCodeSortReq) (*pb.ListBaseCodeSortResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListBaseCodeSortResp{}, nil
}
