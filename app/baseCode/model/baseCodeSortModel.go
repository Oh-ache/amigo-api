package model

import (
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
