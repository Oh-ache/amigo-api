package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeItemLogic {
	return &GetBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeItemLogic) GetBaseCodeItem(in *pb.GetBaseCodeItemReq) (*pb.BaseCodeItemResp, error) {
	var item *model.BaseCodeItem
	var err error

	if in.BaseCodeItemId != 0 {
		item, err = l.svcCtx.BaseCodeItemModel.FindOne(l.ctx, in.BaseCodeItemId)
	} else if in.SortKey != "" && in.Key != "" {
		item, err = l.svcCtx.BaseCodeItemModel.FindOneBySortKeyKey(l.ctx, in.SortKey, in.Key)
	} else {
		return nil, model.ErrInvalidParams
	}

	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrNotFound
		}
		l.Errorf("Failed to get base code item: %v", err)
		return nil, err
	}

	var resp pb.BaseCodeItemResp
	_ = copier.Copy(&resp, item)

	return &resp, nil
}
