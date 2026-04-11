package model

import (
	"context"
	"fmt"
	"strings"

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
	// 构建查询条件
	var conditions []string

	// 处理查询条件（等值查询）
	if search.SortKey != "" {
		conditions = append(conditions, "`sort_key` = '"+search.SortKey+"'")
	}
	if search.Key != "" {
		conditions = append(conditions, "`key` = '"+search.Key+"'")
	}
	if search.ForeignId != "" {
		conditions = append(conditions, "`foreign_id` = '"+search.ForeignId+"'")
	}
	if search.IsDelete != 0 {
		conditions = append(conditions, "`is_delete` = "+fmt.Sprintf("%d", search.IsDelete))
	}
	if search.Content != "" {
		conditions = append(conditions, "`content` like '%"+search.Content+"%'")
	}
	if search.Content1 != "" {
		conditions = append(conditions, "`content1` like '%"+search.Content1+"%'")
	}
	if search.Content2 != "" {
		conditions = append(conditions, "`content2` like '%"+search.Content2+"%'")
	}
	if search.Content3 != "" {
		conditions = append(conditions, "`content3` like '%"+search.Content3+"%'")
	}
	if search.Content4 != "" {
		conditions = append(conditions, "`content4` like '%"+search.Content4+"%'")
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
	query := fmt.Sprintf("select %s from %s %s order by `base_code_item_id` desc %s", baseCodeItemRows, m.table, quertWhere, pageSql)

	// 执行查询
	var list []*BaseCodeItem
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
