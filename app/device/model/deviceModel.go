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

var _ DeviceModel = (*customDeviceModel)(nil)

type (
	// DeviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDeviceModel.
	DeviceModel interface {
		deviceModel
		CheckDuplicate(ctx context.Context, data *Device) (bool, error)
		List(ctx context.Context, search *DeviceSearch) ([]*Device, int64, error)
	}

	DeviceSearch struct {
		Name       string
		UserId     uint64
		MacAddress string
		InternalIp string
		IsRunning  int64
		IsDelete   int64
		Page       int64
		PageSize   int64
	}

	customDeviceModel struct {
		*defaultDeviceModel
	}
)

// NewDeviceModel returns a model for the database table.
func NewDeviceModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) DeviceModel {
	return &customDeviceModel{
		defaultDeviceModel: newDeviceModel(conn, c, opts...),
	}
}

// CheckDuplicate 检查mac_address是否重复（区分新增和更新）
func (m *customDeviceModel) CheckDuplicate(ctx context.Context, data *Device) (bool, error) {
	// 首先通过 macAddress 查找是否存在对应的记录
	existing, err := m.FindOneByMacAddress(ctx, data.MacAddress)
	if err == sqlc.ErrNotFound {
		// 没有找到，说明不重复
		return false, nil
	}
	if err != nil {
		// 查询出错
		return false, err
	}

	// 找到了对应的 macAddress 记录
	if data.DeviceId == 0 {
		// 主键id为0，表示新增数据，只要macAddress存在就认为是重复
		return true, nil
	}

	// 主键id不为0，表示更新数据，需要判断找到的记录是否是自己
	if existing.DeviceId != data.DeviceId {
		// 找到了其他记录拥有相同的macAddress，认为是重复
		return true, nil
	}

	// 找到的记录就是自己，不是重复
	return false, nil
}

// List 方法根据搜索条件查询 Device 列表（按主键id降序）
func (m *customDeviceModel) List(ctx context.Context, search *DeviceSearch) ([]*Device, int64, error) {
	var conditions []string

	if search.Name != "" {
		conditions = append(conditions, "`name` = '"+search.Name+"'")
	}
	if search.UserId != 0 {
		conditions = append(conditions, "`user_id` = "+fmt.Sprintf("%d", search.UserId))
	}
	if search.MacAddress != "" {
		conditions = append(conditions, "`mac_address` = '"+search.MacAddress+"'")
	}
	if search.InternalIp != "" {
		conditions = append(conditions, "`internal_ip` = '"+search.InternalIp+"'")
	}
	if search.IsRunning != 0 {
		conditions = append(conditions, "`is_running` = "+fmt.Sprintf("%d", search.IsRunning))
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
	query := fmt.Sprintf("select %s from %s %s order by `device_id` desc %s", deviceRows, m.table, queryWhere, pageSql)

	var list []*Device
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
