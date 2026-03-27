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
	var pbReq pb.AddBaseCodeReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BaseCodeRpc.AddBaseCode(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
