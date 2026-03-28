package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PngToElnkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPngToElnkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PngToElnkLogic {
	return &PngToElnkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PngToElnkLogic) PngToElnk(in *pb.PngToElnkReq) (*pb.PngToElnkResp, error) {
	// todo: add your logic here and delete this line

	return &pb.PngToElnkResp{}, nil
}
