package ratelimiter

import (
	ratelimiterservice "go-auth-backend-api/internal/service/rateLimiterService"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware(rateLimitKey string, countLimit int64) gin.HandlerFunc {
	return func(c *gin.Context) {

		clientIP := c.ClientIP()

		err := ratelimiterservice.RateLimiterService(clientIP, rateLimitKey, countLimit)

		if err != nil {

			if err.Error() == "Rate Limit exceeded try after 1 minute" {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "try again in 1 minute"})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials"})
			c.Abort()
			return
		}

		c.Next()
	}
}
