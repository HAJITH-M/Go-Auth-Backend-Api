package authhandler

import (
	authservice "go-auth-backend-api/internal/service/authService"
	"go-auth-backend-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChangePasswordHandler(c *gin.Context) {
	var req ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	err := authservice.ChangePasswordService(authservice.ChangePasswordInput{
		Email:       req.Email,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		if err.Error() == "old password doesn't match" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if err.Error() == "New password should be different from old password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Old and New Password are same"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func ForgotPasswordHandler(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	if err := authservice.ForgotPasswordService(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func VerifyForgotPasswordHandler(c *gin.Context) {
	var req VerifyForgotPasswordOtpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	err := authservice.VerifyForgotPasswordOtp(req.Email, req.Otp)
	if err != nil {
		if err.Error() == "Invalid OTP" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP Verification Success"})
}

func ForgotPasswordUpdateHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	err := authservice.ForgotPasswordUpdateService(req.Email, req.Password)
	if err != nil {
		if err.Error() == "new password should not be same as old password" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
