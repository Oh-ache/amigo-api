package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListBaseCodeReq) (resp *types.ListBaseCodeResp, err error) {
	var pbReq pb.ListBaseCodeReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	pbResp, err := l.svcCtx.BaseCodeRpc.ListBaseCode(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	var apiResp types.ListBaseCodeResp
	if err := copier.Copy(&apiResp, pbResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}
