package svc

import (
	"amigo-api/app/ai/model"
	"amigo-api/app/ai/rpc/internal/config"
	"amigo-api/app/baseCode/rpc/basecode"
	"amigo-api/app/sdk/rpc/sdk"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	AiTaskModel  model.AiTaskModel
	BaseCodeRpc  basecode.BaseCode
	SdkRpcClient sdk.Sdk
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config:       c,
		AiTaskModel:  model.NewAiTaskModel(sqlConn, c.Cache),
		BaseCodeRpc:  basecode.NewBaseCode(zrpc.MustNewClient(c.BaseCodeRpcConf)),
		SdkRpcClient: sdk.NewSdk(zrpc.MustNewClient(c.SdkRpcConf)),
	}
}

type AiTask struct {
	Id           int64
	UserId       int64
	TaskId       string
	TaskType     string
	Prompt       string
	RequestInfo  string
	ResponseInfo string
	ResultUrl    string
	Status       int
	ErrorMsg     string
	CreatedAt    int64
	UpdatedAt    int64
}