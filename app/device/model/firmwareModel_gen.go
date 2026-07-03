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
	firmwareFieldNames          = builder.RawFieldNames(&Firmware{})
	firmwareRows                = strings.Join(firmwareFieldNames, ",")
	firmwareRowsExpectAutoSet   = strings.Join(stringx.Remove(firmwareFieldNames, "`firmware_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	firmwareRowsWithPlaceHolder = strings.Join(stringx.Remove(firmwareFieldNames, "`firmware_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheAmigoFirmwareFirmwareIdPrefix = "cache:amigo:firmware:firmwareId:"
)

type (
	firmwareModel interface {
		Insert(ctx context.Context, data *Firmware) (sql.Result, error)
		FindOne(ctx context.Context, firmwareId uint64) (*Firmware, error)
		Update(ctx context.Context, data *Firmware) error
		Delete(ctx context.Context, firmwareId uint64) error
	}

	defaultFirmwareModel struct {
		sqlc.CachedConn
		table string
	}

	Firmware struct {
		FirmwareId uint64 `db:"firmware_id"`
		Name       string `db:"name"`
		Version    string `db:"version"`
		DeviceType string `db:"device_type"`
		FileUrl    string `db:"file_url"`
		FileSize   int64  `db:"file_size"`
		Md5        string `db:"md5"`
		Changelog  string `db:"changelog"`
		IsForce    int64  `db:"is_force"`
		IsDelete   int64  `db:"is_delete"`
		CreateTime uint64 `db:"create_time"`
		UpdateTime uint64 `db:"update_time"`
	}
)

func newFirmwareModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultFirmwareModel {
	return &defaultFirmwareModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`firmware`",
	}
}

func (m *defaultFirmwareModel) Delete(ctx context.Context, firmwareId uint64) error {
	firmwareFirmwareIdKey := fmt.Sprintf("%s%v", cacheAmigoFirmwareFirmwareIdPrefix, firmwareId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `firmware_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, firmwareId)
	}, firmwareFirmwareIdKey)
	return err
}

func (m *defaultFirmwareModel) FindOne(ctx context.Context, firmwareId uint64) (*Firmware, error) {
	firmwareFirmwareIdKey := fmt.Sprintf("%s%v", cacheAmigoFirmwareFirmwareIdPrefix, firmwareId)
	var resp Firmware
	err := m.QueryRowCtx(ctx, &resp, firmwareFirmwareIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `firmware_id` = ? limit 1", firmwareRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, firmwareId)
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

func (m *defaultFirmwareModel) Insert(ctx context.Context, data *Firmware) (sql.Result, error) {
	firmwareFirmwareIdKey := fmt.Sprintf("%s%v", cacheAmigoFirmwareFirmwareIdPrefix, data.FirmwareId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (%s)", m.table, firmwareRowsExpectAutoSet, strings.Repeat("?,", len(stringx.Remove(firmwareFieldNames, "`firmware_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"))-1)+"?")
		return conn.ExecCtx(ctx, query, data.Name, data.Version, data.DeviceType, data.FileUrl, data.FileSize, data.Md5, data.Changelog, data.IsForce, data.IsDelete, data.CreateTime, data.UpdateTime)
	}, firmwareFirmwareIdKey)
	return ret, err
}

func (m *defaultFirmwareModel) Update(ctx context.Context, data *Firmware) error {
	firmwareFirmwareIdKey := fmt.Sprintf("%s%v", cacheAmigoFirmwareFirmwareIdPrefix, data.FirmwareId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `firmware_id` = ?", m.table, firmwareRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Name, data.Version, data.DeviceType, data.FileUrl, data.FileSize, data.Md5, data.Changelog, data.IsForce, data.IsDelete, data.CreateTime, data.UpdateTime, data.FirmwareId)
	}, firmwareFirmwareIdKey)
	return err
}

func (m *defaultFirmwareModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheAmigoFirmwareFirmwareIdPrefix, primary)
}

func (m *defaultFirmwareModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `firmware_id` = ? limit 1", firmwareRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultFirmwareModel) tableName() string {
	return m.table
}
