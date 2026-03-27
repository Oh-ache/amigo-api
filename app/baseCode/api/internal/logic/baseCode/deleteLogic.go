package baseCode

import (
	"context"

	"amigo-api/app/baseCode/api/internal/svc"
	"amigo-api/app/baseCode/api/internal/types"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteBaseCodeReq) (resp *types.EmptyResp, err error) {
	var pbReq pb.DeleteBaseCodeReq
	if err := copier.Copy(&pbReq, req); err != nil {
		return nil, err
	}

	_, err = l.svcCtx.BaseCodeRpc.DeleteBaseCode(l.ctx, &pbReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
