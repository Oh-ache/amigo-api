package logic

import (
	"context"

	"amigo-api/app/baseCode/model"
	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBaseCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBaseCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBaseCodeLogic {
	return &DeleteBaseCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBaseCodeLogic) DeleteBaseCode(in *pb.DeleteBaseCodeReq) (*pb.DeleteBaseCodeResp, error) {
	// 先尝试根据主键id查询数据是否存在
	if in.BaseCodeId == 0 && in.SortKey != "" && in.Key != "" {
		// 主键id不存在，但有sort_key和key，根据sort_key和key查询并获取主键id
		if code, err := l.svcCtx.BaseCodeModel.FindOneBySortKeyKey(l.ctx, in.SortKey, in.Key); err == nil {
			in.BaseCodeId = code.BaseCodeId
		} else if err != model.ErrNotFound {
			l.Errorf("Failed to find BaseCode by sort_key and key: %v", err)
			return &pb.DeleteBaseCodeResp{Success: false}, err
		}
	}

	// 检查主键id是否存在
	if in.BaseCodeId == 0 {
		return &pb.DeleteBaseCodeResp{Success: false}, model.ErrNotFound
	}

	// 根据主键id删除数据
	if err := l.svcCtx.BaseCodeModel.Delete(l.ctx, in.BaseCodeId); err != nil {
		if err == model.ErrNotFound {
			return &pb.DeleteBaseCodeResp{Success: false}, model.ErrNotFound
		}
		l.Errorf("Failed to delete BaseCode by id %d: %v", in.BaseCodeId, err)
		return &pb.DeleteBaseCodeResp{Success: false}, err
	}

	return &pb.DeleteBaseCodeResp{Success: true}, nil
}
