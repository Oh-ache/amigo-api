package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBaseCodeSortLogic {
	return &AddBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBaseCodeSortLogic) AddBaseCodeSort(in *pb.AddBaseCodeSortReq) (*pb.BaseCodeSortResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeSortResp{}, nil
}
