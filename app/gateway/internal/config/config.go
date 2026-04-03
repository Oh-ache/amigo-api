package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Upstreams Upstreams `json:",optional"`
	Auth      Auth      `json:",optional"`
}

type Auth struct {
	AccessSecret string `json:",optional"`
	AccessExpire int64  `json:",optional"`
}

type Upstreams struct {
	User     string `json:",optional"`
	Device   string `json:",optional"`
	Sdk      string `json:",optional"`
	BaseCode string `json:",optional"`
}
