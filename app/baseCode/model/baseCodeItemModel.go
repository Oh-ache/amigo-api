package model

import (
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
