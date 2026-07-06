package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkOrderReplyModel = (*customWorkOrderReplyModel)(nil)

type (
	WorkOrderReplyModel interface {
		workOrderReplyModel
		ListByWorkOrderId(ctx context.Context, workOrderId uint64) ([]*WorkOrderReply, error)
	}

	customWorkOrderReplyModel struct {
		*defaultWorkOrderReplyModel
	}
)

func NewWorkOrderReplyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WorkOrderReplyModel {
	return &customWorkOrderReplyModel{
		defaultWorkOrderReplyModel: newWorkOrderReplyModel(conn, c, opts...),
	}
}

func (m *customWorkOrderReplyModel) ListByWorkOrderId(ctx context.Context, workOrderId uint64) ([]*WorkOrderReply, error) {
	query := fmt.Sprintf("select %s from %s where `work_order_id` = ? and `is_delete` = 2 order by `create_time` asc", workOrderReplyRows, m.table)
	var list []*WorkOrderReply
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, workOrderId); err != nil {
		return nil, err
	}
	return list, nil
}
