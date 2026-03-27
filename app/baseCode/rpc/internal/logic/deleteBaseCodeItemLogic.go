package logic

import (
	"context"

	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/svc"
	"amigo-api/common/pb"

	jsoniter "github.com/json-iterator/go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
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
	// 从上下文中获取tracer
	tracer := trace.TracerFromContext(l.ctx)
	// 创建自定义span
	ctx, span := tracer.Start(l.ctx, "开始删除")
	// 设置span属性

	fast := jsoniter.ConfigFastest
	bytes2, _ := fast.Marshal(in)
	span.SetAttributes(
		attribute.String("deleteItem.param", string(bytes2)),
	)
	defer span.End()

	// 先尝试根据主键id查询数据是否存在
	if in.BaseCodeItemId == 0 && in.SortKey != "" && in.Key != "" {
		// 主键id不存在，但有sort_key和key，根据sort_key和key查询并获取主键id
		if item, err := l.svcCtx.BaseCodeItemModel.FindOneBySortKeyKey(ctx, in.SortKey, in.Key); err == nil {
			in.BaseCodeItemId = item.BaseCodeItemId
		} else if err != model.ErrNotFound {
			l.Errorf("Failed to find BaseCodeItem by sort_key and key: %v", err)
			return &pb.DeleteBaseCodeItemResp{Success: false}, err
		}
	}

	// 检查主键id是否存在
	if in.BaseCodeItemId == 0 {
		return &pb.DeleteBaseCodeItemResp{Success: false}, model.ErrNotFound
	}

	// 根据主键id删除数据
	if err := l.svcCtx.BaseCodeItemModel.Delete(ctx, in.BaseCodeItemId); err != nil {
		if err == model.ErrNotFound {
			return &pb.DeleteBaseCodeItemResp{Success: false}, model.ErrNotFound
		}
		l.Errorf("Failed to delete BaseCodeItem by id %d: %v", in.BaseCodeItemId, err)
		return &pb.DeleteBaseCodeItemResp{Success: false}, err
	}

	// 删除成功
	span.SetAttributes(
		attribute.String("deleteItem.success", "ok"),
	)

	return &pb.DeleteBaseCodeItemResp{Success: true}, nil
}
