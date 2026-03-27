package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.AddBaseCodeReq) (resp *types.EmptyResp, err error) {
	resp = &types.EmptyResp{}
	param := &pb.AddBaseCodeReq{}

	copier.Copy(param, req)
	if _, err := l.svcCtx.BaseCodeRpcClient.AddBaseCode(l.ctx, param); err != nil {
		return nil, err
	}

	return resp, nil
}
