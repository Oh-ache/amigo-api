package model

import (
	"context"
	"fmt"

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
	// 构建查询条件（参数化查询，防止 SQL 注入）
	var conditions []utils.SQLCondition

	if search.SortKey != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`sort_key` = ?", Args: []any{search.SortKey}})
	}
	if search.SortName != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`sort_name` = ?", Args: []any{search.SortName}})
	}
	if search.IsDelete != 0 {
		conditions = append(conditions, utils.SQLCondition{Clause: "`is_delete` = ?", Args: []any{search.IsDelete}})
	}

	queryWhere, whereArgs := utils.BuildWhereClause(conditions)

	// 获取总条数（不算分页的数据）
	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, queryWhere)

	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, whereArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// 构建查询语句（按主键id降序）
	pageClause, pageArgs := utils.DelSQLPage(search.Page, search.PageSize)
	query := fmt.Sprintf("select %s from %s %s order by `base_code_sort_id` desc %s", baseCodeSortRows, m.table, queryWhere, pageClause)

	// 合并所有参数
	args := append(whereArgs, pageArgs...)

	// 执行查询
	var list []*BaseCodeSort
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
