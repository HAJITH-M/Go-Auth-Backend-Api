package redis

import (
	"context"
	"fmt"
	redisconfig "go-auth-backend-api/internal/config/redisConfig"
	"go-auth-backend-api/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

func StoreOtp(email string, otp string, expirationTime time.Duration) error {
	ctx := context.Background()

	otpHashed, err := utils.GeneratePasswordWithHash(otp)
	if err != nil {
		return err
	}

	otpCoolDown := "otp_cooldown:" + email

	err = redisconfig.RDB.Set(ctx, otpCoolDown, "1", 30*time.Second).Err()
	if err != nil {
		return err
	}

	err = redisconfig.RDB.Set(ctx, email, otpHashed, expirationTime).Err()
	if err != nil {
		return err

	}

	fmt.Print(redisconfig.RDB.Get(ctx, email).Result())

	return nil
}

func GetOtp(OtpEmail string) (string, error) {

	ctx := context.Background()

	val, err := redisconfig.RDB.Get(ctx, OtpEmail).Result()

	if err == redis.Nil {
		// OTP not found (expired or not created)
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func DeleteOtp(key string) error {

	ctx := context.Background()

	return redisconfig.RDB.Del(ctx, key).Err()
}

func CheckOtpCooldown(email string) (bool, error) {
	ctx := context.Background()

	// cooldown key must match the key used in StoreOtp:
	// otpCoolDown := "otp_cooldown:" + emailParam
	// where emailParam is currently passed as "otp:"+email from ForgotPasswordService.
	// So here we need to check "otp_cooldown:otp:<email>".
	key := "otp_cooldown:otp:" + email

	exists, err := redisconfig.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func IncrOtpAttempts(email string, otpTTL time.Duration) (int64, error) {
	ctx := context.Background()
	key := "otp:attempts:" + email

	attempts, err := redisconfig.RDB.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment attempts: %w", err)
	}
	if attempts == 1 {
		redisconfig.RDB.Expire(ctx, key, otpTTL)
	}

	return attempts, nil
}

func DeleteOtpAttempts(email string) {
	ctx := context.Background()
	redisconfig.RDB.Del(ctx, "otp:attempts:"+email)
}

func RateLimiterRedis(clientIP string, rateLimitKey string, rateLimitTTL time.Duration) (int64, error) {
	ctx := context.Background()
	key := rateLimitKey + clientIP

	attempts, err := redisconfig.RDB.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment attempts: %w", err)
	}

	if attempts == 1 {
		redisconfig.RDB.Expire(ctx, key, rateLimitTTL)
	}

	return attempts, nil
}
