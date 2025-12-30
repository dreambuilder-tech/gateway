package member

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/currency"
	"wallet/common-lib/rpcx/wallet_rpcx"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type AssetsResp struct {
	Currency currency.Code   `json:"currency"`
	Free     decimal.Decimal `json:"free"`
	CNY      decimal.Decimal `json:"cny"`
}

func Assets(c *gin.Context) {
	memberID := auth.MemberID(c)
	resp, err := wallet_rpcx.GetWallet(c.Request.Context(), memberID, currency.Coin)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, &AssetsResp{
		Currency: currency.Coin,
		Free:     resp.Free,
		CNY:      resp.Free,
	})
}
