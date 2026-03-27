package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBaseCodeSortLogic {
	return &DeleteBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBaseCodeSortLogic) DeleteBaseCodeSort(in *pb.DeleteBaseCodeSortReq) (*pb.DeleteBaseCodeSortResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DeleteBaseCodeSortResp{}, nil
}
