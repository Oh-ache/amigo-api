package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// Redis configuration
	Redis HostPort

	// MQueue configuration
	MQueue MQueueConfig
}

type HostPort struct {
	Host string
	Port int
}

type MQueueConfig struct {
	RedisHost   string
	RedisPort   int
	RedisPass   string
	RedisDB     int
	ServerName  string
	Queues      map[string]int
	Concurrency int
	SyncTimeout time.Duration
	RetryDelay  time.Duration
	MaxRetry    int
	Timeout     time.Duration
	DeadQueue   string
	MonitorPort int
}

func (c *MQueueConfig) GetRedisOptions() *RedisOptions {
	return &RedisOptions{
		Addr:     c.RedisHost,
		Port:     c.RedisPort,
		Password: c.RedisPass,
		DB:       c.RedisDB,
	}
}

type RedisOptions struct {
	Addr     string
	Port     int
	Password string
	DB       int
}
