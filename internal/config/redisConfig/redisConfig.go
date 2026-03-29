package redisconfig

import (
	"go-auth-backend-api/internal/config/env"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func RedisConfig() *redis.Client {
	RDB = redis.NewClient(&redis.Options{
		Addr:     env.AppEnv.REDIS_ADDR,
		Username: env.AppEnv.REDIS_USERNAME,
		Password: env.AppEnv.REDIS_PASSWORD,
		DB:       env.AppEnv.REDIS_DB,
	})

	return RDB
}
