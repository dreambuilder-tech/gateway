package router

import (
	"gateway/internal/handler/notify"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func notifyRouter(r *gin.RouterGroup) {
	ra := r.Group("/", middleware.Auth())
	{
		ra.POST("/sync", notify.Sync)
		ra.POST("/ack", notify.Ack)

		rac := ra.Group("/center")
		{
			rac.GET("/unread", notify.CenterUnread)
			rac.POST("/list", notify.CenterList)
		}
	}
}
