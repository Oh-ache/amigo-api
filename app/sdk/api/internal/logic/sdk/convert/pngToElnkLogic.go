package convert

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PngToElnkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPngToElnkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PngToElnkLogic {
	return &PngToElnkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PngToElnkLogic) PngToElnk(req *types.PngToElnkReq) (resp *types.PngToElnkResp, err error) {
	// todo: add your logic here and delete this line

	return
}
