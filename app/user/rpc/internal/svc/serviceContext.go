package svc

import (
	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/app/user/model"
	"amigo-api/app/user/rpc/internal/config"
	"amigo-api/common/utils"
	"amigo-api/common/utils/plug/userauth"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	RedisClient    *redis.Redis
	UserModel      model.UserModel
	AdminModel     model.AdminModel
	UserThirdParty model.UserThirdPartyModel
	AdminAuth      *userauth.UserAuthClient
	BaseCodeRpc    basecode.BaseCode
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	conf, _ := utils.ReadFileToString("etc/casbin.conf")
	adminAuthClient, _ := userauth.NewClient(c.DB.DataSource, conf)

	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),

		UserModel:      model.NewUserModel(sqlConn, c.Cache),
		AdminModel:     model.NewAdminModel(sqlConn, c.Cache),
		UserThirdParty: model.NewUserThirdPartyModel(sqlConn, c.Cache),

		AdminAuth: adminAuthClient,

		BaseCodeRpc: basecode.NewBaseCode(zrpc.MustNewClient(c.BaseCodeRpcConf)),
	}
}
