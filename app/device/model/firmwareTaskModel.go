package model

import (
	"context"
	"fmt"
	"strings"

	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FirmwareTaskModel = (*customFirmwareTaskModel)(nil)

type (
	FirmwareTaskModel interface {
		firmwareTaskModel
		List(ctx context.Context, search *FirmwareTaskSearch) ([]*FirmwareTask, int64, error)
	}

	FirmwareTaskSearch struct {
		FirmwareId uint64
		DeviceId   uint64
		Status     int64
		Page       int64
		PageSize   int64
	}

	customFirmwareTaskModel struct {
		*defaultFirmwareTaskModel
	}
)

func NewFirmwareTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FirmwareTaskModel {
	return &customFirmwareTaskModel{
		defaultFirmwareTaskModel: newFirmwareTaskModel(conn, c, opts...),
	}
}

func (m *customFirmwareTaskModel) List(ctx context.Context, search *FirmwareTaskSearch) ([]*FirmwareTask, int64, error) {
	var conditions []string

	if search.FirmwareId != 0 {
		conditions = append(conditions, "`firmware_id` = "+fmt.Sprintf("%d", search.FirmwareId))
	}
	if search.DeviceId != 0 {
		conditions = append(conditions, "`device_id` = "+fmt.Sprintf("%d", search.DeviceId))
	}
	if search.Status != 0 {
		conditions = append(conditions, "`status` = "+fmt.Sprintf("%d", search.Status))
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
	query := fmt.Sprintf("select %s from %s %s order by `firmware_task_id` desc %s", firmwareTaskRows, m.table, queryWhere, pageSql)

	var list []*FirmwareTask
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
