package model

import (
	"context"
	"fmt"
	"strings"

	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FirmwareModel = (*customFirmwareModel)(nil)

type (
	FirmwareModel interface {
		firmwareModel
		CheckDuplicate(ctx context.Context, data *Firmware) (bool, error)
		List(ctx context.Context, search *FirmwareSearch) ([]*Firmware, int64, error)
	}

	FirmwareSearch struct {
		DeviceType string
		Version    string
		Name       string
		IsForce    int64
		IsDelete   int64
		Page       int64
		PageSize   int64
	}

	customFirmwareModel struct {
		*defaultFirmwareModel
	}
)

func NewFirmwareModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FirmwareModel {
	return &customFirmwareModel{
		defaultFirmwareModel: newFirmwareModel(conn, c, opts...),
	}
}

func (m *customFirmwareModel) CheckDuplicate(ctx context.Context, data *Firmware) (bool, error) {
	query := fmt.Sprintf("select %s from %s where `device_type` = ? and `version` = ? and `is_delete` = 2 limit 1", firmwareRows, m.table)

	var existing Firmware
	err := m.QueryRowNoCacheCtx(ctx, &existing, query, data.DeviceType, data.Version)
	if err == sqlc.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if data.FirmwareId == 0 {
		return true, nil
	}

	if existing.FirmwareId != data.FirmwareId {
		return true, nil
	}

	return false, nil
}

func (m *customFirmwareModel) List(ctx context.Context, search *FirmwareSearch) ([]*Firmware, int64, error) {
	var conditions []string

	if search.DeviceType != "" {
		conditions = append(conditions, "`device_type` = '"+search.DeviceType+"'")
	}
	if search.Version != "" {
		conditions = append(conditions, "`version` = '"+search.Version+"'")
	}
	if search.Name != "" {
		conditions = append(conditions, "`name` like '%"+search.Name+"%'")
	}
	if search.IsForce != 0 {
		conditions = append(conditions, "`is_force` = "+fmt.Sprintf("%d", search.IsForce))
	}
	if search.IsDelete != 0 {
		conditions = append(conditions, "`is_delete` = "+fmt.Sprintf("%d", search.IsDelete))
	}

	queryWhere := ""
	if len(conditions) > 0 {
		queryWhere = " where " + strings.Join(conditions, " and ")
	}

	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, queryWhere)

	var total int64
	if m.QueryRowNoCacheCtx(ctx, &total, countQuery) != nil {
		return nil, 0, fmt.Errorf("failed to get total count")
	}

	pageSql := utils.DelSQLPage(search.Page, search.PageSize)
	query := fmt.Sprintf("select %s from %s %s order by `firmware_id` desc %s", firmwareRows, m.table, queryWhere, pageSql)

	var list []*Firmware
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
