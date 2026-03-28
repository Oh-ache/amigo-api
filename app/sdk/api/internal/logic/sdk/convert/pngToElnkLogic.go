package convert

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"
	"amigo-api/app/sdk/rpc/sdk"

	"github.com/jinzhu/copier"
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
	resp = &types.PngToElnkResp{}
	param := &sdk.PngToElnkReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.SdkRpcClient.PngToElnk(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
