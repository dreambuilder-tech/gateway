package transfer

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/currency"
	"wallet/common-lib/consts/ledger"
	"wallet/common-lib/dto/res_dto/state"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/merchant_rpcx"
	"wallet/common-lib/rpcx/trade_rpcx"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type ResolveResp struct {
	state.BaseResp
	Display    string          `json:"display,omitempty"`
	PayeeID    int64           `json:"payee_id,omitempty"`
	PayeeCode  string          `json:"payee_code,omitempty"`
	Amount     decimal.Decimal `json:"amount,omitempty"`
	ExpireAtMS int64           `json:"expire_at_ms,omitempty"`
}

func Resolve(c *gin.Context) {
	tk := c.Query("t")
	if tk == "" {
		app.InvalidParams(c, "empty token")
		return
	}
	resp, err := merchant_rpcx.InternalResolveToken(c.Request.Context(), tk)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	r := &ResolveResp{
		BaseResp:   state.BaseResp{State: resp.State},
		Display:    resp.Display,
		PayeeID:    resp.PayeeID,
		PayeeCode:  resp.PayeeCode,
		Amount:     resp.Amount,
		ExpireAtMS: resp.ExpireAtMS,
	}
	app.Result(c, r)
}

type ConfirmReq struct {
	Token string `json:"token"`
	PIN   string `json:"pin"`
}

func Confirm(c *gin.Context) {
	req := new(ConfirmReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.Token == "" || req.PIN == "" {
		app.InvalidParams(c, "empty parameters")
		return
	}
	TK, err := merchant_rpcx.InternalResolveToken(c.Request.Context(), req.Token)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	if !TK.OK() {
		app.Result(c, TK)
		return
	}
	resp, err := trade_rpcx.MerchantTopup(c.Request.Context(), &contracts.MemberTransferReq{
		PIN:   req.PIN,
		Payee: TK.PayeeID,
		TransferReqCommon: contracts.TransferReqCommon{
			OrderNo:   TK.TopupID,
			From:      auth.MemberID(c),
			PayeeCode: TK.PayeeCode,
			Reason:    ledger.Merchant_Member_Topup,
			Currency:  currency.Coin,
			Amount:    TK.Amount,
			Desc:      "商户玩家充值",
		},
	})
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
