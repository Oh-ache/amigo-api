package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	rest.RestConf
	Redis RedisConf
	Queue QueueConf
}

type RedisConf struct {
	redis.RedisConf
	DB int
}

type QueueConf struct {
	Prefix         string
	DefaultQueue   string
	RetryQueue     string
	DelayQueue     string
	DeadLetterQueue string
	MaxRetry       int
	Concurrency    int
	PollInterval   int
}