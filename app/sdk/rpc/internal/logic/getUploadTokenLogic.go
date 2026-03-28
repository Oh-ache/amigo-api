package logic

import (
	"context"

	"amigo-api/app/sdk/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadTokenLogic {
	return &GetUploadTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUploadTokenLogic) GetUploadToken(in *pb.GetUploadTokenReq) (*pb.GetUploadTokenResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUploadTokenResp{}, nil
}
