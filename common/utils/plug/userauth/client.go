// Package  u  provides ...
package userauth

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type UserAuthClient struct {
	Enforcer *casbin.Enforcer
}

func NewClient(db, modelStr string) (*UserAuthClient, error) {
	conn, _ := xorm.NewEngine("mysql", db)
	adapter, err := xormadapter.NewAdapterByEngine(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %v", err)
	}

	m, err := model.NewModelFromString(modelStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %v", err)
	}

	// 加载RBAC带有域的策略模型
	enforcer, err := casbin.NewEnforcer(m, adapter)
	enforcer.EnableLog(false)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %v", err)
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %v", err)
	}

	return &UserAuthClient{Enforcer: enforcer}, nil
}
