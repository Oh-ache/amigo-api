package model

import (
	"context"
	"fmt"
	"strings"

	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BaseCodeSortModel = (*customBaseCodeSortModel)(nil)

type (
	// BaseCodeSortModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBaseCodeSortModel.
	BaseCodeSortModel interface {
		baseCodeSortModel
	}

	customBaseCodeSortModel struct {
		*defaultBaseCodeSortModel
	}
)

// NewBaseCodeSortModel returns a model for the database table.
func NewBaseCodeSortModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BaseCodeSortModel {
	return &customBaseCodeSortModel{
		defaultBaseCodeSortModel: newBaseCodeSortModel(conn, c, opts...),
	}
}

// List 方法根据搜索条件查询 BaseCodeSort 列表（按主键id降序）
func (m *customBaseCodeSortModel) List(ctx context.Context, search *BaseCodeSortSearch) ([]*BaseCodeSort, int64, error) {
	// 构建查询条件
	var conditions []string

	// 处理前三个查询条件（等值查询）
	if search.SortKey != "" {
		conditions = append(conditions, "`sort_key` = '"+search.SortKey+"'")
	}
	if search.SortName != "" {
		conditions = append(conditions, "`sort_name` = '"+search.SortName+"'")
	}
	if search.IsDelete != 0 {
		conditions = append(conditions, "`is_delete` = "+fmt.Sprintf("%d", search.IsDelete))
	}

	quertWhere := ""
	if len(conditions) > 0 {
		quertWhere = " where " + strings.Join(conditions, " and ")
	}

	// 获取总条数（不算分页的数据）
	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, quertWhere)

	var total int64
	if m.QueryRowNoCacheCtx(ctx, &total, countQuery) != nil {
		return nil, 0, fmt.Errorf("failed to get total count")
	}

	// 构建查询语句（按主键id降序）
	pageSql := utils.DelSQLPage(search.Page, search.PageSize)
	query := fmt.Sprintf("select %s from %s %s order by `base_code_sort_id` desc %s", baseCodeSortRows, m.table, quertWhere, pageSql)

	// 执行查询
	var list []*BaseCodeSort
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
