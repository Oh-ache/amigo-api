package oss

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"
	"amigo-api/app/sdk/rpc/sdk"

	"github.com/jinzhu/copier"
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
	resp = &types.UploadTokenResp{}
	param := &sdk.GetUploadTokenReq{}

	copier.Copy(param, req)
	rpcResp, err := l.svcCtx.SdkRpcClient.GetUploadToken(l.ctx, param)
	if err != nil {
		return nil, err
	}

	copier.Copy(resp, rpcResp)
	return resp, nil
}
