package svc

import (
	"amigo-api/app/device/model"
	"amigo-api/app/device/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	RedisClient  *redis.Redis
	DeviceModel  model.DeviceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		DeviceModel: model.NewDeviceModel(sqlConn, c.Cache),
	}
}
