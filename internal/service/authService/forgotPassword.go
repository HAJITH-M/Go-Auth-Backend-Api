package authservice

import (
	"errors"
	"go-auth-backend-api/internal/repository"
	"go-auth-backend-api/pkg/mailer"
	"go-auth-backend-api/pkg/redis"
	"go-auth-backend-api/pkg/utils"
	"time"
)

func ForgotPasswordService(email string) error {
	cooldown, err := redis.CheckOtpCooldown(email)
	if err != nil {
		return err
	}

	if cooldown {
		return errors.New("please wait 30 seconds before requesting another OTP")
	}

	otpStr, err := utils.GenerateOTP()
	if err != nil {
		return err
	}

	ctx := "otp:" + email
	err = redis.StoreOtp(ctx, otpStr, 5*time.Minute)
	if err != nil {
		return errors.New("Failed to store OTP")
	}

	err = repository.UpdateUserAccountStatusRepo(email, "pending_otp")
	if err != nil {
		return err
	}

	go func() {
		if sendErr := mailer.SendOtpEmail(email, otpStr); sendErr != nil {
			// intentionally ignored to keep request non-blocking
		}
	}()

	return nil
}

func VerifyForgotPasswordOtp(email string, otp string) error {
	attempts, err := redis.IncrOtpAttempts(email, 5*time.Minute)
	if err != nil {
		return err
	}
	if attempts > 3 {
		redis.DeleteOtp("otp:" + email)
		redis.DeleteOtpAttempts(email)
		return errors.New("too many attempts, please request a new OTP")
	}

	getOtp, err := redis.GetOtp("otp:" + email)
	if err != nil {
		return err
	}

	err = utils.CompareHashedPassword(getOtp, otp)
	if err != nil {
		return errors.New("Invalid OTP")
	}

	err = redis.DeleteOtp("otp:" + email)
	if err != nil {
		return err
	}

	redis.DeleteOtpAttempts(email)
	return nil
}
