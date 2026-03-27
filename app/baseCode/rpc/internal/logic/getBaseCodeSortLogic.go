package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeSortLogic {
	return &GetBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeSortLogic) GetBaseCodeSort(in *pb.GetBaseCodeSortReq) (*pb.BaseCodeSortResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeSortResp{}, nil
}
