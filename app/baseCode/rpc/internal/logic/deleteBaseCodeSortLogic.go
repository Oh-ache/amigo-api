package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBaseCodeSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBaseCodeSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBaseCodeSortLogic {
	return &DeleteBaseCodeSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBaseCodeSortLogic) DeleteBaseCodeSort(in *pb.DeleteBaseCodeSortReq) (*pb.DeleteBaseCodeSortResp, error) {
	// 先尝试根据主键id查询数据是否存在
	if in.BaseCodeSortId == 0 && in.SortKey != "" {
		// 主键id不存在，但有sort_key，根据sort_key查询并获取主键id
		if sort, err := l.svcCtx.BaseCodeSortModel.FindOneBySortKey(l.ctx, in.SortKey); err == nil {
			in.BaseCodeSortId = sort.BaseCodeSortId
		} else if err != model.ErrNotFound {
			l.Errorf("Failed to find BaseCodeSort by sort_key: %v", err)
			return &pb.DeleteBaseCodeSortResp{Success: false}, err
		}
	}

	// 检查主键id是否存在
	if in.BaseCodeSortId == 0 {
		return &pb.DeleteBaseCodeSortResp{Success: false}, model.ErrNotFound
	}

	// 根据主键id删除数据
	if err := l.svcCtx.BaseCodeSortModel.Delete(l.ctx, in.BaseCodeSortId); err != nil {
		if err == model.ErrNotFound {
			return &pb.DeleteBaseCodeSortResp{Success: false}, model.ErrNotFound
		}
		l.Errorf("Failed to delete BaseCodeSort by id %d: %v", in.BaseCodeSortId, err)
		return &pb.DeleteBaseCodeSortResp{Success: false}, err
	}

	return &pb.DeleteBaseCodeSortResp{Success: true}, nil
}
