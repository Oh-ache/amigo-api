package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	workOrderFieldNames          = builder.RawFieldNames(&WorkOrder{})
	workOrderRows                = strings.Join(workOrderFieldNames, ",")
	workOrderRowsExpectAutoSet   = strings.Join(stringx.Remove(workOrderFieldNames, "`work_order_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	workOrderRowsWithPlaceHolder = strings.Join(stringx.Remove(workOrderFieldNames, "`work_order_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheAmigoWorkOrderWorkOrderIdPrefix = "cache:amigo:workOrder:workOrderId:"
)

type (
	workOrderModel interface {
		Insert(ctx context.Context, data *WorkOrder) (sql.Result, error)
		FindOne(ctx context.Context, workOrderId uint64) (*WorkOrder, error)
		Update(ctx context.Context, data *WorkOrder) error
		Delete(ctx context.Context, workOrderId uint64) error
	}

	defaultWorkOrderModel struct {
		sqlc.CachedConn
		table string
	}

	WorkOrder struct {
		WorkOrderId uint64 `db:"work_order_id"`
		DeviceId    uint64 `db:"device_id"`
		UserId      uint64 `db:"user_id"`
		Title       string `db:"title"`
		Content     string `db:"content"`
		Images      string `db:"images"`
		Category    int64  `db:"category"`
		Status      int64  `db:"status"`
		IsDelete    int64  `db:"is_delete"`
		CreateTime  uint64 `db:"create_time"`
		UpdateTime  uint64 `db:"update_time"`
	}
)

func newWorkOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultWorkOrderModel {
	return &defaultWorkOrderModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`work_order`",
	}
}

func (m *defaultWorkOrderModel) Delete(ctx context.Context, workOrderId uint64) error {
	key := fmt.Sprintf("%s%v", cacheAmigoWorkOrderWorkOrderIdPrefix, workOrderId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `work_order_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, workOrderId)
	}, key)
	return err
}

func (m *defaultWorkOrderModel) FindOne(ctx context.Context, workOrderId uint64) (*WorkOrder, error) {
	key := fmt.Sprintf("%s%v", cacheAmigoWorkOrderWorkOrderIdPrefix, workOrderId)
	var resp WorkOrder
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `work_order_id` = ? limit 1", workOrderRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, workOrderId)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWorkOrderModel) Insert(ctx context.Context, data *WorkOrder) (sql.Result, error) {
	key := fmt.Sprintf("%s%v", cacheAmigoWorkOrderWorkOrderIdPrefix, data.WorkOrderId)
	removedFields := stringx.Remove(workOrderFieldNames, "`work_order_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`")
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (%s)", m.table, workOrderRowsExpectAutoSet, strings.Repeat("?,", len(removedFields)-1)+"?")
		return conn.ExecCtx(ctx, query, data.DeviceId, data.UserId, data.Title, data.Content, data.Images, data.Category, data.Status, data.IsDelete, data.CreateTime, data.UpdateTime)
	}, key)
	return ret, err
}

func (m *defaultWorkOrderModel) Update(ctx context.Context, data *WorkOrder) error {
	key := fmt.Sprintf("%s%v", cacheAmigoWorkOrderWorkOrderIdPrefix, data.WorkOrderId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `work_order_id` = ?", m.table, workOrderRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.DeviceId, data.UserId, data.Title, data.Content, data.Images, data.Category, data.Status, data.IsDelete, data.CreateTime, data.UpdateTime, data.WorkOrderId)
	}, key)
	return err
}

func (m *defaultWorkOrderModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheAmigoWorkOrderWorkOrderIdPrefix, primary)
}

func (m *defaultWorkOrderModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `work_order_id` = ? limit 1", workOrderRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}
