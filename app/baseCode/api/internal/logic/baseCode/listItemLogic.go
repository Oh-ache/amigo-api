package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
