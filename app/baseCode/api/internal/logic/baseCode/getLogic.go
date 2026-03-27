package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogic) Get(req *types.GetBaseCodeReq) (resp *types.GetBaseCodeResp, err error) {
	var pbReq pb.GetBaseCodeReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	pbResp, err := l.svcCtx.BaseCodeRpc.GetBaseCode(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	var apiResp types.GetBaseCodeResp
	if err := copier.Copy(&apiResp, pbResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}
