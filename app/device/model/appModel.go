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

var _ AppModel = (*customAppModel)(nil)

type (
	// AppModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAppModel.
	AppModel interface {
		appModel
		CheckDuplicate(ctx context.Context, data *App) (bool, error)
		List(ctx context.Context, search *AppSearch) ([]*App, int64, error)
	}

	AppSearch struct {
		Name       string
		AppKey     string
		Category   string
		SensorType string
		IsDelete   int64
		Page       int64
		PageSize   int64
	}

	customAppModel struct {
		*defaultAppModel
	}
)

// NewAppModel returns a model for the database table.
func NewAppModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AppModel {
	return &customAppModel{
		defaultAppModel: newAppModel(conn, c, opts...),
	}
}

// CheckDuplicate 检查 app_key 是否重复（区分新增和更新）
func (m *customAppModel) CheckDuplicate(ctx context.Context, data *App) (bool, error) {
	existing, err := m.FindOneByAppKey(ctx, data.AppKey)
	if err == sqlc.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if data.AppId == 0 {
		return true, nil
	}
	if existing.AppId != data.AppId {
		return true, nil
	}
	return false, nil
}

// List 方法根据搜索条件查询 App 列表（按主键id降序）
func (m *customAppModel) List(ctx context.Context, search *AppSearch) ([]*App, int64, error) {
	var conditions []string
	if search.Name != "" {
		conditions = append(conditions, "`name` like '%"+search.Name+"%'")
	}
	if search.AppKey != "" {
		conditions = append(conditions, "`app_key` = '"+search.AppKey+"'")
	}
	if search.Category != "" {
		conditions = append(conditions, "`category` = '"+search.Category+"'")
	}
	if search.SensorType != "" {
		conditions = append(conditions, "`sensor_type` = '"+search.SensorType+"'")
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
	query := fmt.Sprintf("select %s from %s %s order by `app_id` desc %s", appRows, m.table, queryWhere, pageSql)
	var list []*App
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
