package oss

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadUrlLogic {
	return &UploadUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadUrlLogic) UploadUrl(req *types.UploadUrlReq) (resp *types.UploadUrlResp, err error) {
	// todo: add your logic here and delete this line

	return
}
