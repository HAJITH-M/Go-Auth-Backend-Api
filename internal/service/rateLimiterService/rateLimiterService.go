package ratelimiterservice

import (
	autherrors "go-auth-backend-api/internal/errors"
	"go-auth-backend-api/pkg/redis"
	"time"
)

func RateLimiterService(clientIP string, rateLimitKey string, countLimit int64) error {

	attempts, err := redis.RateLimiterRedis(clientIP, rateLimitKey, time.Minute)
	if err != nil {
		return autherrors.ErrRateLimiterRedis
	}

	if attempts >= countLimit {
		return autherrors.ErrRateLimitExceeded
	}

	return nil
}
