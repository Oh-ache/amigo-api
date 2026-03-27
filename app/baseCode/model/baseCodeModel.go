package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BaseCodeModel = (*customBaseCodeModel)(nil)

type (
	// BaseCodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBaseCodeModel.
	BaseCodeModel interface {
		baseCodeModel
	}

	customBaseCodeModel struct {
		*defaultBaseCodeModel
	}
)

// NewBaseCodeModel returns a model for the database table.
func NewBaseCodeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) BaseCodeModel {
	return &customBaseCodeModel{
		defaultBaseCodeModel: newBaseCodeModel(conn, c, opts...),
	}
}
