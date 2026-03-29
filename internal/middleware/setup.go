package middleware

import (
	"go-auth-backend-api/internal/middleware/cors"
	middleware "go-auth-backend-api/internal/middleware/logger"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(cors.CorsMiddleware())
}
