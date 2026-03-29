package router

import (
	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			authRoutes(v1)
		}
	}
}
