package oss

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadTokenLogic {
	return &UploadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadTokenLogic) UploadToken(req *types.UploadTokenReq) (resp *types.UploadTokenResp, err error) {
	// todo: add your logic here and delete this line

	return
}
