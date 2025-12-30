package router

import (
	"gateway/internal/handler/transfer"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func transferRouter(r *gin.RouterGroup) {
	auth := r.Group("/", middleware.Auth())
	ts := auth.Group("/transfer")
	{
		ts.GET("/preflight", transfer.Preflight)
		ts.GET("/resolve", transfer.ResolvePayeeCode)
		ts.POST("/confirm", transfer.ConfirmTransfer)
	}
	tco := auth.Group("/checkout")
	{
		tco.GET("/resolve", transfer.Resolve)
		tco.POST("/confirm", transfer.Confirm)
	}
}
