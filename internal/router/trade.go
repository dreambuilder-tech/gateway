package router

import (
	"gateway/internal/handler/trade"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func tradeRouter(r *gin.RouterGroup) {
	t := r.Group("/")
	{
		t.GET("/price", trade.GetPrice)
	}
	auth := t.Group("/", middleware.Auth())

	tc := auth.Group("/c2c")
	{
		tca := tc.Group("/ads")
		{
			tca.GET("/list", trade.GetC2cAdsList)
			tca.GET("/detail", trade.GetC2cAdsDetail)
			tca.POST("/create", trade.CreateC2cAds)
			tca.POST("/pause", trade.PauseC2cAds)
			tca.POST("/publish", trade.PublishC2cAds)
			tca.POST("/cancel", trade.CancelC2cAds)
		}
		tco := tc.Group("/order")
		{
			tco.GET("/detail", trade.GetC2cOrderDetail)
			tco.POST("/create", trade.CreateC2cOrder)
			tco.POST("/agree", trade.AgreeC2cOrder)
			tco.POST("/cancel", trade.CancelC2cOrder)
			tco.POST("/pay", trade.PayC2cOrder)
			tco.POST("/pay/cert", trade.UpdateC2cPayCert)
			tco.POST("/send", trade.SendC2cOrder)
			tco.POST("/unusual", trade.UnusualC2cOrder)
			tco.GET("/history", trade.GetC2cOrderHistory)
		}
		tcp := tc.Group("/price")
		{
			tcp.GET("/list", trade.GetC2cPriceList)
		}
		tcm := tc.Group("/my")
		{
			tcm.GET("/ads/list", trade.GetMyC2cAdsList)
			tcm.GET("/order/list", trade.GetMyC2cOrderList)
		}
	}
}
