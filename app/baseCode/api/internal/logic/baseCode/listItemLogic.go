package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListItemLogic {
	return &ListItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListItemLogic) ListItem(req *types.ListBaseCodeItemReq) (resp *types.ListBaseCodeItemResp, err error) {
	var pbReq pb.ListBaseCodeItemReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	pbResp, err := l.svcCtx.BaseCodeRpc.ListBaseCodeItem(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	var apiResp types.ListBaseCodeItemResp
	if err := copier.Copy(&apiResp, pbResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}
