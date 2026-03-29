package authhandler

import (
	"go-auth-backend-api/internal/repository"
	authservice "go-auth-backend-api/internal/service/authService"
	"go-auth-backend-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	// Check if user exists and has only a Google provider — block password login.
	userByEmail, err := repository.GetUserByEmailRepo(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials"})
		return
	}
	if userByEmail != nil {
		googleMethod, err := repository.GetAuthenticationMethodByUserAndProviderRepo(userByEmail.ID, "google")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials"})
			return
		}
		passwordMethod, err := repository.GetAuthenticationMethodByUserAndProviderRepo(userByEmail.ID, "password")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials"})
			return
		}

		// Google-only user trying to login with password.
		// Safe to show Google message here because:
		// - We already confirmed the email exists
		// - We already confirmed it has a Google provider
		// - The user is not authenticated yet but this helps UX, not attackers
		if googleMethod != nil && passwordMethod == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":    "invalid credentials",
				"hint":     "This account was created using Google. Please sign in with Google.",
				"provider": "google",
			})
			return
		}
	}

	res, err := authservice.LoginService(authservice.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		IPAddress: c.ClientIP(),
		Device:    c.Request.UserAgent(),
	})

	if err != nil {
		if err.Error() == "Email Not Verified" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "account is not active" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "Change password Required" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.SetCookie(
		"refresh_token",
		res.RefreshToken,
		7*24*3600,
		"/api/v1/auth/refresh",
		"",
		false,
		true,
	)
	c.JSON(http.StatusOK, gin.H{
		"message":      "login successful",
		"access_token": res.AccessToken,
		"user": gin.H{
			"user_id":      res.UserId,
			"email":        res.Email,
			"display_name": res.DisplayName,
		},
	})
}
