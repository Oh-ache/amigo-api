package logic

import (
	"context"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBaseCodeSortLogic {
	return &UpdateBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBaseCodeSortLogic) UpdateBaseCodeSort(in *pb.BaseCodeSortResp) (*pb.BaseCodeSortResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BaseCodeSortResp{}, nil
}
