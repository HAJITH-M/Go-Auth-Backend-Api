package authhandler

import (
	authservice "go-auth-backend-api/internal/service/authService"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	result, err := authservice.RefreshTokenService(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refresh_token",
		result.RefreshToken,
		7*24*3600,
		"/api/v1/auth/refresh",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": result.AccessToken,
	})
}

func LogoutHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		authservice.LogoutService(refreshToken)
	}

	c.SetCookie("refresh_token", "", -1, "/api/v1/auth/refresh", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func UpdateVerificationEmailHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "token not found"})
		return
	}

	if err := authservice.UpdateVerificationEmailService(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "email verified successfully",
	})
}
