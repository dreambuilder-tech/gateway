package trade

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/currency"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/trade_rpcx"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func GetC2cAdsList(c *gin.Context) {
	req := contracts.C2CAdsListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.PageArgs.Init()
	if req.Side == 0 {
		req.Side = 1
	}
	if req.Currency == "" {
		req.Currency = currency.Coin
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.GetC2cAdsList(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.ResultPage(c, resp.List, resp.Total)
}

func GetC2cAdsDetail(c *gin.Context) {
	req := contracts.C2CAdsDetailReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.GetC2cAdsDetail(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func CreateC2cAds(c *gin.Context) {
	req := contracts.C2CCreateAdsReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if !req.Price.IsPositive() {
		app.InvalidParams(c, "price is not positive")
		return
	}
	if !req.Total.IsPositive() {
		app.InvalidParams(c, "total is not positive")
		return
	}
	if !req.MinCount.IsPositive() {
		req.MinCount = decimal.NewFromFloat(100)
	}
	if !req.MaxCount.IsPositive() {
		req.MaxCount = decimal.NewFromFloat(99999999)
	}
	if req.MinCount.LessThan(decimal.NewFromFloat(100)) {
		app.InvalidParams(c, "min_count is less than 100")
		return
	}
	if req.MaxCount.LessThan(req.MinCount) {
		app.InvalidParams(c, "max_count is less than min_count")
		return
	}
	if len(req.PayMethods) == 0 {
		app.InvalidParams(c, "pay_methods is empty")
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CCreateAds(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func PauseC2cAds(c *gin.Context) {
	req := contracts.C2CPauseAdsReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.AdsID <= 0 {
		app.InvalidParams(c, "ads_id is less than 0")
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CPauseAds(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func PublishC2cAds(c *gin.Context) {
	req := contracts.C2CPublishAdsReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.AdsID <= 0 {
		app.InvalidParams(c, "ads_id is less than 0")
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CPublishAds(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func CancelC2cAds(c *gin.Context) {
	req := contracts.C2CCancelAdsReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.AdsID <= 0 {
		app.InvalidParams(c, "ads_id is less than 0")
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CCancelAds(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func GetC2cOrderDetail(c *gin.Context) {
	req := contracts.C2COrderDetailReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.ID <= 0 && req.OrderID == "" {
		app.InvalidParams(c, "id or order_id is empty")
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.GetC2cOrderDetail(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func CreateC2cOrder(c *gin.Context) {
	req := contracts.C2CCreateOrderReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.Total.LessThanOrEqual(decimal.NewFromFloat(0)) {
		app.InvalidParams(c, "total is less than or equal 0")
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CCreateOrder(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func AgreeC2cOrder(c *gin.Context) {
	req := contracts.C2CAgreeOrderReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CAgreeOrder(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func CancelC2cOrder(c *gin.Context) {
	req := contracts.C2CCancelOrderReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CCancelOrder(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func PayC2cOrder(c *gin.Context) {
	req := contracts.C2CPayOrderReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CPayOrder(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func UpdateC2cPayCert(c *gin.Context) {
	req := contracts.C2CPayCertReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CPayCert(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

func SendC2cOrder(c *gin.Context) {
	req := contracts.C2CSendOrderReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CSendOrder(c.Request.Context(), &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.Result(c, resp)
}

func UnusualC2cOrder(c *gin.Context) {
	req := contracts.C2CUnusualOrderReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.C2CUnusualOrder(c, &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
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

	resp, err := trade_rpcx.GetC2cOrderHistory(c, &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.ResultPage(c, resp.List, resp.Total)
}

func GetC2cPriceList(c *gin.Context) {
	req := contracts.C2CEmptyReq{}
	resp, err := trade_rpcx.GetC2cPriceList(c, &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.ResultPage(c, resp.List, resp.Total)
}

func GetMyC2cAdsList(c *gin.Context) {
	req := contracts.MyC2cAdsListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.PageArgs.Init()
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.GetMyC2cAdsList(c, &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.ResultPage(c, resp.List, resp.Total)
}

func GetMyC2cOrderList(c *gin.Context) {
	req := contracts.MyC2cOrderListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.PageArgs.Init()
	req.MemberID = auth.MemberID(c)

	resp, err := trade_rpcx.GetMyC2cOrderList(c, &req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}

	app.ResultPage(c, resp.List, resp.Total)
}
