package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserThirdPartyModel = (*customUserThirdPartyModel)(nil)

type (
	// UserThirdPartyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserThirdPartyModel.
	UserThirdPartyModel interface {
		userThirdPartyModel
	}

	customUserThirdPartyModel struct {
		*defaultUserThirdPartyModel
	}
)

// NewUserThirdPartyModel returns a model for the database table.
func NewUserThirdPartyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserThirdPartyModel {
	return &customUserThirdPartyModel{
		defaultUserThirdPartyModel: newUserThirdPartyModel(conn, c, opts...),
	}
}
