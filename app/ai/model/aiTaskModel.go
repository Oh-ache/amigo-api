package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AiTaskModel = (*customAiTaskModel)(nil)

type (
	AiTaskModel interface {
		aiTaskModel
	}

	customAiTaskModel struct {
		*defaultAiTaskModel
	}
)

func NewAiTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AiTaskModel {
	return &customAiTaskModel{
		defaultAiTaskModel: newAiTaskModel(conn, c, opts...),
	}
}
