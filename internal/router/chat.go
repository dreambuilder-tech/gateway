package router

import (
	"gateway/internal/handler/chat"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func chatRouter(r *gin.RouterGroup) {
	im := r.Group("/", middleware.Auth())
	{
		im.GET("/ticket", chat.Ticket)
		im.POST("/init", chat.Init)
		im.POST("/history", chat.History)
	}
}
