package router

import (
	authhandler "go-auth-backend-api/internal/handler/authHandler"
	oauthhandler "go-auth-backend-api/internal/handler/oAuthHandler"
	ratelimiter "go-auth-backend-api/internal/middleware/rateLimiter"

	"github.com/gin-gonic/gin"
)

func authRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.Use(ratelimiter.RateLimiterMiddleware("loginRateLimit:", 6))

	{
		auth.GET("/google", oauthhandler.GoogleLogin)
		auth.GET("/google/callback", oauthhandler.GoogleCallback)
		auth.POST("/google/refresh", oauthhandler.RefreshAccessToken)
		auth.POST("/register", authhandler.RegisterHandler)
		auth.POST("/login", authhandler.LoginHandler)
		auth.POST("/refresh", authhandler.RefreshHandler)
		auth.POST("/logout", authhandler.LogoutHandler)
		auth.GET("/verification-email", authhandler.UpdateVerificationEmailHandler)
		auth.POST("/change-password", authhandler.ChangePasswordHandler)
		auth.POST("/forgot-password", authhandler.ForgotPasswordHandler)
		auth.POST("/forgot-password-verify", authhandler.VerifyForgotPasswordHandler)
		auth.POST("/forgot-password-update", authhandler.ForgotPasswordUpdateHandler)
	}
}
