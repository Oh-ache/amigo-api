package model

import (
	"context"
	"fmt"

	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BaseCodeItemModel = (*customBaseCodeItemModel)(nil)

type (
	// BaseCodeItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBaseCodeItemModel.
	BaseCodeItemModel interface {
		baseCodeItemModel
	}

	customBaseCodeItemModel struct {
		*defaultBaseCodeItemModel
	}
)

// NewBaseCodeItemModel returns a model for the database table.
func NewBaseCodeItemModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BaseCodeItemModel {
	return &customBaseCodeItemModel{
		defaultBaseCodeItemModel: newBaseCodeItemModel(conn, c, opts...),
	}
}

// List 方法根据搜索条件查询 BaseCodeItem 列表（按主键id降序）
func (m *customBaseCodeItemModel) List(ctx context.Context, search *BaseCodeItemSearch) ([]*BaseCodeItem, int64, error) {
	// 构建查询条件（参数化查询，防止 SQL 注入）
	var conditions []utils.SQLCondition

	if search.SortKey != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`sort_key` = ?", Args: []any{search.SortKey}})
	}
	if search.Key != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`key` = ?", Args: []any{search.Key}})
	}
	if search.ForeignId != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`foreign_id` = ?", Args: []any{search.ForeignId}})
	}
	if search.IsDelete != 0 {
		conditions = append(conditions, utils.SQLCondition{Clause: "`is_delete` = ?", Args: []any{search.IsDelete}})
	}
	if search.Content != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`content` like ?", Args: []any{fmt.Sprintf("%%%s%%", search.Content)}})
	}
	if search.Content1 != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`content1` like ?", Args: []any{fmt.Sprintf("%%%s%%", search.Content1)}})
	}
	if search.Content2 != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`content2` like ?", Args: []any{fmt.Sprintf("%%%s%%", search.Content2)}})
	}
	if search.Content3 != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`content3` like ?", Args: []any{fmt.Sprintf("%%%s%%", search.Content3)}})
	}
	if search.Content4 != "" {
		conditions = append(conditions, utils.SQLCondition{Clause: "`content4` like ?", Args: []any{fmt.Sprintf("%%%s%%", search.Content4)}})
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
	query := fmt.Sprintf("select %s from %s %s order by `base_code_item_id` desc %s", baseCodeItemRows, m.table, queryWhere, pageClause)

	// 合并所有参数
	args := append(whereArgs, pageArgs...)

	// 执行查询
	var list []*BaseCodeItem
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
