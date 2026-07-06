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
	workOrderReplyFieldNames          = builder.RawFieldNames(&WorkOrderReply{})
	workOrderReplyRows                = strings.Join(workOrderReplyFieldNames, ",")
	workOrderReplyRowsExpectAutoSet   = strings.Join(stringx.Remove(workOrderReplyFieldNames, "`work_order_reply_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	workOrderReplyRowsWithPlaceHolder = strings.Join(stringx.Remove(workOrderReplyFieldNames, "`work_order_reply_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheAmigoWorkOrderReplyWorkOrderReplyIdPrefix = "cache:amigo:workOrderReply:workOrderReplyId:"
)

type (
	workOrderReplyModel interface {
		Insert(ctx context.Context, data *WorkOrderReply) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*WorkOrderReply, error)
		Update(ctx context.Context, data *WorkOrderReply) error
	}

	defaultWorkOrderReplyModel struct {
		sqlc.CachedConn
		table string
	}

	WorkOrderReply struct {
		WorkOrderReplyId uint64 `db:"work_order_reply_id"`
		WorkOrderId      uint64 `db:"work_order_id"`
		AdminId          uint64 `db:"admin_id"`
		Content          string `db:"content"`
		Images           string `db:"images"`
		IsDelete         int64  `db:"is_delete"`
		CreateTime       uint64 `db:"create_time"`
		UpdateTime       uint64 `db:"update_time"`
	}
)

func newWorkOrderReplyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultWorkOrderReplyModel {
	return &defaultWorkOrderReplyModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`work_order_reply`",
	}
}

func (m *defaultWorkOrderReplyModel) FindOne(ctx context.Context, id uint64) (*WorkOrderReply, error) {
	key := fmt.Sprintf("%s%v", cacheAmigoWorkOrderReplyWorkOrderReplyIdPrefix, id)
	var resp WorkOrderReply
	err := m.QueryRowCtx(ctx, &resp, key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `work_order_reply_id` = ? limit 1", workOrderReplyRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
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

func (m *defaultWorkOrderReplyModel) Insert(ctx context.Context, data *WorkOrderReply) (sql.Result, error) {
	key := fmt.Sprintf("%s%v", cacheAmigoWorkOrderReplyWorkOrderReplyIdPrefix, data.WorkOrderReplyId)
	removedFields := stringx.Remove(workOrderReplyFieldNames, "`work_order_reply_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`")
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (%s)", m.table, workOrderReplyRowsExpectAutoSet, strings.Repeat("?,", len(removedFields)-1)+"?")
		return conn.ExecCtx(ctx, query, data.WorkOrderId, data.AdminId, data.Content, data.Images, data.IsDelete, data.CreateTime, data.UpdateTime)
	}, key)
	return ret, err
}

func (m *defaultWorkOrderReplyModel) Update(ctx context.Context, data *WorkOrderReply) error {
	key := fmt.Sprintf("%s%v", cacheAmigoWorkOrderReplyWorkOrderReplyIdPrefix, data.WorkOrderReplyId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `work_order_reply_id` = ?", m.table, workOrderReplyRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.WorkOrderId, data.AdminId, data.Content, data.Images, data.IsDelete, data.CreateTime, data.UpdateTime, data.WorkOrderReplyId)
	}, key)
	return err
}

func (m *defaultWorkOrderReplyModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheAmigoWorkOrderReplyWorkOrderReplyIdPrefix, primary)
}

func (m *defaultWorkOrderReplyModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `work_order_reply_id` = ? limit 1", workOrderReplyRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}
