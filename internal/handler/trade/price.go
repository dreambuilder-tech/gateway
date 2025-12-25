package trade

import (
	"gateway/internal/app"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/trade_rpcx"

	"github.com/gin-gonic/gin"
)

func GetPrice(c *gin.Context) {
	req := contracts.C2CPriceReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}

	resp, err := trade_rpcx.C2CPrice(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.SuccessData(c, resp)
}
