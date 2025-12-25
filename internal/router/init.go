package router

import (
	"gateway/internal/common/rule"
	// "gateway/internal/handler/chat"

	"gateway/internal/handler/member"
	// "gateway/internal/handler/trade"

	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) {
	engine.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	wsRouter(engine.Group("/ws"))

	r := engine.Group("/api/v1")
	memberRouter(r.Group("/member"))
	// tradeRouter(r.Group("/trade"))
	// chatRouter(r.Group("/chat"))
}

func memberRouter(r *gin.RouterGroup) {
	m := r.Group("/")
	{
		m.POST("/send-sms", withLimit(rule.SendSMS), member.SendSms)
		m.POST("/register", withLimit(rule.Register), member.Register)
		m.POST("/login", withLimit(rule.Login), member.Login)
		m.POST("/reset-pwd", member.ResetPwd)
		m.POST("/change-pwd", member.ChangePwd)
	}
}

// // func tradeRouter(r *gin.RouterGroup) {
// // 	t := r.Group("/")
// // 	{
// // 		t.GET("/price", trade.GetPrice)
// // 	}

// // 	tc := r.Group("/c2c", middleware.Auth())
// // 	{
// // 		tc.GET("/list", trade.GetC2cAdsList)
// // 		tc.GET("/history", trade.GetC2cOrderHistory)
// // 	}
// // }

// func chatRouter(r *gin.RouterGroup) {
// 	im := r.Group("/", middleware.Auth())
// 	{
// 		im.GET("/ticket", chat.Ticket)
// 		im.POST("/init", chat.Init)
// 		im.POST("/history", chat.History)
// 	}
// }
