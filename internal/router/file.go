package router

import (
	"gateway/internal/handler/files"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func filesRouter(r *gin.RouterGroup) {
	auth := r.Group("/", middleware.Auth())
	{
		auth.POST("/upload", files.Upload)
	}
}
