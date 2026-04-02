package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Upstreams Upstreams `json:",optional"`
}

type Upstreams struct {
	UserAdmin    string `json:",optional"`
	Device       string `json:",optional"`
	Sdk          string `json:",optional"`
	BaseCode     string `json:",optional"`
}
