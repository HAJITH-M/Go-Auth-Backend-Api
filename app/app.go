package app

import (
	"go-auth-backend-api/internal/config/env"
	redisconfig "go-auth-backend-api/internal/config/redisConfig"
	"go-auth-backend-api/internal/middleware"
	"go-auth-backend-api/internal/router"
	"go-auth-backend-api/pkg/database"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	bootstrapOnce sync.Once
	handler       http.Handler
	bootstrapErr  error
)

func Handler() (http.Handler, error) {
	bootstrapOnce.Do(func() {
		gin.SetMode(gin.ReleaseMode)

		r := gin.New()
		_ = r.SetTrustedProxies(nil)

		env.Load()
		redisconfig.RedisConfig()

		if err := database.Connect(); err != nil {
			bootstrapErr = err
			return
		}

		log.Println("Database connected successfully")
		middleware.Setup(r)
		router.SetUpRouter(r)

		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "hello I'm Go-Auth-Backend-Api Backend"})
		})

		handler = r
	})

	return handler, bootstrapErr
}
