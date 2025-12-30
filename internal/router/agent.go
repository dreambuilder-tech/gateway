package router

import (
	"gateway/internal/handler/agent"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func agentRouter(r *gin.RouterGroup) {
	aa := r.Group("/", middleware.Auth())
	{
		aa.GET("/preflight", agent.Preflight)
		aa.POST("/apply", agent.Apply)
		aa.POST("/refund", agent.Refund)
	}
}
