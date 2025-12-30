package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) {
	engine.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	wsRouter(engine.Group("/ws"))

	r := engine.Group("/api/v1")

	filesRouter(r.Group("/files"))
	memberRouter(r.Group("/member"))
	tradeRouter(r.Group("/trade"))
	chatRouter(r.Group("/chat"))
	transferRouter(r)
	agentRouter(r.Group("/agent"))
	notifyRouter(r.Group("/notify"))
}
