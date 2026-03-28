package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadUrlLogic {
	return &UploadUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadUrlLogic) UploadUrl(in *pb.UploadUrlReq) (*pb.UploadUrlResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UploadUrlResp{}, nil
}
