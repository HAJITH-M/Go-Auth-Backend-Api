package ratelimiterservice

import (
	"errors"
	"go-auth-backend-api/pkg/redis"
	"time"
)

func RateLimiterService(clientIP string, rateLimitKey string, countLimit int64) error {

	attempts, err := redis.RateLimiterRedis(clientIP, rateLimitKey, time.Minute)
	if err != nil {
		return errors.New("failed to generate token")
	}

	if attempts >= countLimit {
		return errors.New("Rate Limit exceeded try after 1 minute")
	}

	return nil
}
