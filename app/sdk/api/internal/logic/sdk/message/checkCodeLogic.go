package message

import (
	"context"

	"amigo-api/app/sdk/api/internal/svc"
	"amigo-api/app/sdk/api/internal/types"
	"amigo-api/app/sdk/rpc/sdk"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckCodeLogic {
	return &CheckCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckCodeLogic) CheckCode(req *types.CheckCodeReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	param := &sdk.CheckCodeReq{}

	copier.Copy(param, req)
	if _, err := l.svcCtx.SdkRpcClient.CheckCode(l.ctx, param); err != nil {
		return nil, err
	}

	return resp, nil
}
