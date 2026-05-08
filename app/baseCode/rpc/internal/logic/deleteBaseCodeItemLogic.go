package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBaseCodeItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBaseCodeItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBaseCodeItemLogic {
	return &DeleteBaseCodeItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBaseCodeItemLogic) DeleteBaseCodeItem(in *pb.DeleteBaseCodeItemReq) (*pb.DeleteBaseCodeItemResp, error) {
	resp := &pb.DeleteBaseCodeItemResp{Success: false}
	// 先尝试根据主键id查询数据是否存在
	if in.BaseCodeItemId == 0 && in.SortKey != "" && in.Key != "" {
		// 主键id不存在，但有sort_key和key，根据sort_key和key查询并获取主键id
		if item, err := l.svcCtx.BaseCodeItemModel.FindOneBySortKeyKey(l.ctx, in.SortKey, in.Key); err == nil {
			in.BaseCodeItemId = item.BaseCodeItemId
		} else if err != model.ErrNotFound {
			l.Errorf("Failed to find BaseCodeItem by sort_key and key: %v", err)
			return resp, err
		}
	}

	// 检查主键id是否存在
	if in.BaseCodeItemId == 0 {
		return resp, model.ErrNotFound
	}

	// 根据主键id删除数据
	if err := l.svcCtx.BaseCodeItemModel.Delete(l.ctx, in.BaseCodeItemId); err != nil {
		if err == model.ErrNotFound {
			return resp, model.ErrNotFound
		}
		l.Errorf("Failed to delete BaseCodeItem by id %d: %v", in.BaseCodeItemId, err)
		return resp, err
	}

	resp.Success = true
	return resp, nil
}
