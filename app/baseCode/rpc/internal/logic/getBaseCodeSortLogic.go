package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeSortLogic {
	return &GetBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeSortLogic) GetBaseCodeSort(in *pb.GetBaseCodeSortReq) (*pb.BaseCodeSortResp, error) {
	var sort *model.BaseCodeSort
	var err error

	// 优先根据主键id查询
	if in.BaseCodeSortId != 0 {
		if sort, err = l.svcCtx.BaseCodeSortModel.FindOne(l.ctx, in.BaseCodeSortId); err == nil {
			return l.buildResponse(l.ctx, sort)
		} else if err != model.ErrNotFound {
			return nil, err
		}
	}

	// 其次根据sort_key查询
	if in.SortKey != "" {
		if sort, err = l.svcCtx.BaseCodeSortModel.FindOneBySortKey(l.ctx, in.SortKey); err == nil {
			return l.buildResponse(l.ctx, sort)
		} else if err != model.ErrNotFound {
			return nil, err
		}
		return nil, model.ErrNotFound
	}

	// 未提供查询条件
	return nil, model.ErrNotFound
}

// 统一构建响应对象
func (l *GetBaseCodeSortLogic) buildResponse(ctx context.Context, sort *model.BaseCodeSort) (*pb.BaseCodeSortResp, error) {
	var resp pb.BaseCodeSortResp
	if err := copier.Copy(&resp, sort); err != nil {
		return nil, err
	}
	return &resp, nil
}
