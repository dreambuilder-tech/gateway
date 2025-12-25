package trade

import (
	"gateway/internal/app"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/trade_rpcx"

	"github.com/gin-gonic/gin"
)

func GetC2cAdsList(c *gin.Context) {
	req := contracts.C2CListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.Coin == "" {
		app.InvalidParams(c, "currency is empty")
		return
	}
	if req.Limit <= 0 {
		app.InvalidParams(c, "limit is unavailable")
		return
	}
	if req.Side == 0 {
		app.InvalidParams(c, "side is empty")
		return
	}
	if req.CurrencyMax.LessThan(req.CurrencyMin) {
		app.InvalidParams(c, "currency max less than currency min")
		return
	}

	resp, err := trade_rpcx.C2CList(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.SuccessData(c, resp)
}

func GetC2cOrderHistory(c *gin.Context) {
	req := contracts.C2COrderHistoryReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.PageArgs.Init()
	req.TimeRange.Init()

	resp, err := trade_rpcx.C2COrderHistory(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.SuccessPage(c, resp.List, resp.Total)
}
