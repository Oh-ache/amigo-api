package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Routes []RouteRule `json:",optional"`
	Auth   Auth        `json:",optional"`
}

type RouteRule struct {
	Prefix   string `json:",optional"`
	Upstream string `json:",optional"`
}

type Auth struct {
	AccessSecret string `json:",optional"`
	AccessExpire int64  `json:",optional"`
}
