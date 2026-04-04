package handler

import "github.com/redis/go-redis/v9"

var (
	RedisClient *redis.Client
)

func InitRedis(redisClient *redis.Client) {
	RedisClient = redisClient
}