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
	firmwareTaskFieldNames          = builder.RawFieldNames(&FirmwareTask{})
	firmwareTaskRows                = strings.Join(firmwareTaskFieldNames, ",")
	firmwareTaskRowsExpectAutoSet   = strings.Join(stringx.Remove(firmwareTaskFieldNames, "`firmware_task_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	firmwareTaskRowsWithPlaceHolder = strings.Join(stringx.Remove(firmwareTaskFieldNames, "`firmware_task_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheAmigoFirmwareTaskFirmwareTaskIdPrefix = "cache:amigo:firmwareTask:firmwareTaskId:"
)

type (
	firmwareTaskModel interface {
		Insert(ctx context.Context, data *FirmwareTask) (sql.Result, error)
		FindOne(ctx context.Context, firmwareTaskId uint64) (*FirmwareTask, error)
		Update(ctx context.Context, data *FirmwareTask) error
		Delete(ctx context.Context, firmwareTaskId uint64) error
	}

	defaultFirmwareTaskModel struct {
		sqlc.CachedConn
		table string
	}

	FirmwareTask struct {
		FirmwareTaskId uint64 `db:"firmware_task_id"`
		FirmwareId     uint64 `db:"firmware_id"`
		DeviceId       uint64 `db:"device_id"`
		Status         int64  `db:"status"`
		Progress       int64  `db:"progress"`
		ErrorMsg       string `db:"error_msg"`
		StartedAt      uint64 `db:"started_at"`
		CompletedAt    uint64 `db:"completed_at"`
		IsDelete       int64  `db:"is_delete"`
		CreateTime     uint64 `db:"create_time"`
		UpdateTime     uint64 `db:"update_time"`
	}
)

func newFirmwareTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultFirmwareTaskModel {
	return &defaultFirmwareTaskModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`firmware_task`",
	}
}

func (m *defaultFirmwareTaskModel) Delete(ctx context.Context, firmwareTaskId uint64) error {
	firmwareTaskIdKey := fmt.Sprintf("%s%v", cacheAmigoFirmwareTaskFirmwareTaskIdPrefix, firmwareTaskId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `firmware_task_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, firmwareTaskId)
	}, firmwareTaskIdKey)
	return err
}

func (m *defaultFirmwareTaskModel) FindOne(ctx context.Context, firmwareTaskId uint64) (*FirmwareTask, error) {
	firmwareTaskIdKey := fmt.Sprintf("%s%v", cacheAmigoFirmwareTaskFirmwareTaskIdPrefix, firmwareTaskId)
	var resp FirmwareTask
	err := m.QueryRowCtx(ctx, &resp, firmwareTaskIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `firmware_task_id` = ? limit 1", firmwareTaskRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, firmwareTaskId)
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

func (m *defaultFirmwareTaskModel) Insert(ctx context.Context, data *FirmwareTask) (sql.Result, error) {
	firmwareTaskIdKey := fmt.Sprintf("%s%v", cacheAmigoFirmwareTaskFirmwareTaskIdPrefix, data.FirmwareTaskId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		removedFields := stringx.Remove(firmwareTaskFieldNames, "`firmware_task_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`")
		query := fmt.Sprintf("insert into %s (%s) values (%s)", m.table, firmwareTaskRowsExpectAutoSet, strings.Repeat("?,", len(removedFields)-1)+"?")
		return conn.ExecCtx(ctx, query, data.FirmwareId, data.DeviceId, data.Status, data.Progress, data.ErrorMsg, data.StartedAt, data.CompletedAt, data.IsDelete, data.CreateTime, data.UpdateTime)
	}, firmwareTaskIdKey)
	return ret, err
}

func (m *defaultFirmwareTaskModel) Update(ctx context.Context, data *FirmwareTask) error {
	firmwareTaskIdKey := fmt.Sprintf("%s%v", cacheAmigoFirmwareTaskFirmwareTaskIdPrefix, data.FirmwareTaskId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `firmware_task_id` = ?", m.table, firmwareTaskRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.FirmwareId, data.DeviceId, data.Status, data.Progress, data.ErrorMsg, data.StartedAt, data.CompletedAt, data.IsDelete, data.CreateTime, data.UpdateTime, data.FirmwareTaskId)
	}, firmwareTaskIdKey)
	return err
}

func (m *defaultFirmwareTaskModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheAmigoFirmwareTaskFirmwareTaskIdPrefix, primary)
}

func (m *defaultFirmwareTaskModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `firmware_task_id` = ? limit 1", firmwareTaskRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultFirmwareTaskModel) tableName() string {
	return m.table
}
