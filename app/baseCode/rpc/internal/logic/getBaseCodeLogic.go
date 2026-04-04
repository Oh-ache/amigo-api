package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBaseCodeLogic {
	return &GetBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBaseCodeLogic) GetBaseCode(in *pb.GetBaseCodeReq) (*pb.BaseCodeResp, error) {
	// 根据主键id查询数据
	if in.BaseCodeId != 0 {
		data, err := l.svcCtx.BaseCodeModel.FindOne(l.ctx, in.BaseCodeId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			l.Errorf("Failed to find base code by id %d: %v", in.BaseCodeId, err)
			return nil, err
		}
		var resp pb.BaseCodeResp
		if err := copier.Copy(&resp, data); err != nil {
			l.Errorf("Failed to copy model to response: %v", err)
			return nil, err
		}
		return &resp, nil
	}

	// 根据sort_key和key查询数据
	if in.SortKey != "" && in.Key != "" {
		data, err := l.svcCtx.BaseCodeModel.FindOneBySortKeyKey(l.ctx, in.SortKey, in.Key)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, model.ErrNotFound
			}
			return nil, err
		}
		var resp pb.BaseCodeResp
		if err := copier.Copy(&resp, data); err != nil {
			return nil, err
		}
		return &resp, nil
	}

	return nil, model.ErrNotFound
}
