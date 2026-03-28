package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCourierLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourierLogic {
	return &GetCourierLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourierLogic) GetCourier(in *pb.GetCourierReq) (*pb.GetCourierResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetCourierResp{}, nil
}
