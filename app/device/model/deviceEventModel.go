package model

import (
	"context"
	"fmt"
	"strings"

	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DeviceEventModel = (*customDeviceEventModel)(nil)

type (
	// DeviceEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDeviceEventModel.
	DeviceEventModel interface {
		deviceEventModel
		List(ctx context.Context, search *DeviceEventSearch) ([]*DeviceEvent, int64, error)
		DeleteByDevice(ctx context.Context, deviceId uint64) error
	}

	DeviceEventSearch struct {
		DeviceId   uint64
		EventType  string
		EventLevel string
		Source     string
		IsDelete   int64
		Page       int64
		PageSize   int64
	}

	customDeviceEventModel struct {
		*defaultDeviceEventModel
	}
)

// NewDeviceEventModel returns a model for the database table.
func NewDeviceEventModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) DeviceEventModel {
	return &customDeviceEventModel{
		defaultDeviceEventModel: newDeviceEventModel(conn, c, opts...),
	}
}

// List 根据搜索条件查询设备事件列表（按主键id降序）
func (m *customDeviceEventModel) List(ctx context.Context, search *DeviceEventSearch) ([]*DeviceEvent, int64, error) {
	var conditions []string
	if search.DeviceId != 0 {
		conditions = append(conditions, "`device_id` = "+fmt.Sprintf("%d", search.DeviceId))
	}
	if search.EventType != "" {
		conditions = append(conditions, "`event_type` = '"+search.EventType+"'")
	}
	if search.EventLevel != "" {
		conditions = append(conditions, "`event_level` = '"+search.EventLevel+"'")
	}
	if search.Source != "" {
		conditions = append(conditions, "`source` = '"+search.Source+"'")
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
	query := fmt.Sprintf("select %s from %s %s order by `device_event_id` desc %s", deviceEventRows, m.table, queryWhere, pageSql)

	var list []*DeviceEvent
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// DeleteByDevice 按设备ID清空该设备的所有事件
func (m *customDeviceEventModel) DeleteByDevice(ctx context.Context, deviceId uint64) error {
	query := fmt.Sprintf("delete from %s where `device_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, deviceId)
	return err
}
