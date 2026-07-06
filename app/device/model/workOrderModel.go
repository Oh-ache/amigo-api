package model

import (
	"context"
	"fmt"
	"strings"

	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkOrderModel = (*customWorkOrderModel)(nil)

type (
	WorkOrderModel interface {
		workOrderModel
		List(ctx context.Context, search *WorkOrderSearch) ([]*WorkOrder, int64, error)
	}

	WorkOrderSearch struct {
		DeviceId uint64
		UserId   uint64
		Status   int64
		Category int64
		Page     int64
		PageSize int64
	}

	customWorkOrderModel struct {
		*defaultWorkOrderModel
	}
)

func NewWorkOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WorkOrderModel {
	return &customWorkOrderModel{
		defaultWorkOrderModel: newWorkOrderModel(conn, c, opts...),
	}
}

func (m *customWorkOrderModel) List(ctx context.Context, search *WorkOrderSearch) ([]*WorkOrder, int64, error) {
	var conditions []string

	if search.DeviceId != 0 {
		conditions = append(conditions, "`device_id` = "+fmt.Sprintf("%d", search.DeviceId))
	}
	if search.UserId != 0 {
		conditions = append(conditions, "`user_id` = "+fmt.Sprintf("%d", search.UserId))
	}
	if search.Status != 0 {
		conditions = append(conditions, "`status` = "+fmt.Sprintf("%d", search.Status))
	}
	if search.Category != 0 {
		conditions = append(conditions, "`category` = "+fmt.Sprintf("%d", search.Category))
	}
	conditions = append(conditions, "`is_delete` = 2")

	queryWhere := " where " + strings.Join(conditions, " and ")

	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, queryWhere)
	var total int64
	if m.QueryRowNoCacheCtx(ctx, &total, countQuery) != nil {
		return nil, 0, fmt.Errorf("failed to get total count")
	}

	pageSql := utils.DelSQLPage(search.Page, search.PageSize)
	query := fmt.Sprintf("select %s from %s %s order by `work_order_id` desc %s", workOrderRows, m.table, queryWhere, pageSql)

	var list []*WorkOrder
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
