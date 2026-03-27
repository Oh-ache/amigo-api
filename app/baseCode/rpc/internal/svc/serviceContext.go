package svc

import (
	"amigo-api/app/baseCode/model"

	"amigo-api/app/baseCode/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	RedisClient       *redis.Redis
	BaseCodeSortModel model.BaseCodeSortModel
	BaseCodeModel     model.BaseCodeModel
	BaseCodeItemModel model.BaseCodeItemModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		BaseCodeSortModel: model.NewBaseCodeSortModel(sqlConn, c.Cache),
		BaseCodeModel:     model.NewBaseCodeModel(sqlConn, c.Cache),
		BaseCodeItemModel: model.NewBaseCodeItemModel(sqlConn, c.Cache),
	}
}
