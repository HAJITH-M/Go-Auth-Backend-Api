package main

import (
	"go-auth-backend-api/internal/config/env"
	redisconfig "go-auth-backend-api/internal/config/redisConfig"
	"go-auth-backend-api/internal/router"

	"go-auth-backend-api/internal/middleware"
	"go-auth-backend-api/pkg/database"

	// "go-auth-backend-api/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 1. Load configs first
	env.Load()
	redisconfig.RedisConfig()

	// 2. Connect DB
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	log.Println("✅ Database connected successfully")

	// 3. Setup middleware (after everything is ready)
	middleware.Setup(r)
	router.SetUpRouter(r)

	// 4. Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello I'm ATS Backend"})
	})

	log.Println("🚀 Server running on :8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
