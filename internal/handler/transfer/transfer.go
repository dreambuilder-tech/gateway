package transfer

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/codex"
	"wallet/common-lib/consts/currency"
	"wallet/common-lib/consts/ledger"
	"wallet/common-lib/consts/member_status"
	"wallet/common-lib/errs"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/rpcx/trade_rpcx"
	"wallet/common-lib/rpcx/wallet_rpcx"
	"wallet/common-lib/utils/radx"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// These structs are used for Swagger documentation alongside the original types
// The original types with custom types are kept for business logic compatibility
// DecimalString represents a decimal value as a string in responses
type DecimalString string
type StatusCode string

type PreflightResp struct {
	MemberStatus member_status.Code `json:"member_status"`
	HasPIN       bool               `json:"has_pin"`
	Free         decimal.Decimal    `json:"free"`
	Frozen       decimal.Decimal    `json:"frozen"`
	OrderNo      string             `json:"order_no"`
}

func Preflight(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		memberID = auth.MemberID(c)
	)
	member, err := member_rpcx.MemberInfo(ctx, memberID)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	w, err := wallet_rpcx.GetWallet(ctx, memberID, currency.Coin)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	resp := &PreflightResp{
		MemberStatus: member.Status,
		HasPIN:       member.PinHash != "",
		Free:         w.Free,
		Frozen:       w.Frozen,
	}
	if member.Role.IsAgent() {
		resp.OrderNo = radx.GenTransferOrderNo(1)
	} else {
		resp.OrderNo = radx.GenTransferOrderNo(0)
	}
	app.Result(c, resp)
}

type ResolvePayeeCodeResp struct {
	PayeeStatus member_status.Code `json:"payee_status"`
	Display     string             `json:"display,omitempty"`
	PayeeID     int64              `json:"payee_id,omitempty"`
	PayeeCode   string             `json:"payee_code,omitempty"`
}

func ResolvePayeeCode(c *gin.Context) {
	payeeCode := c.Query("p")
	if payeeCode == "" {
		app.InvalidParams(c, "empty payee code")
		return
	}
	member, err := member_rpcx.MemberInfoByPayeeCode(c.Request.Context(), payeeCode)
	if err != nil {
		if errs.IsRecordNotFound(err) {
			app.Fail(c, codex.PayeeCodeNotFound, "invalid payee code")
			return
		}
		app.InternalError(c, err.Error())
		return
	}
	resp := &ResolvePayeeCodeResp{
		PayeeStatus: member.Status,
		Display:     member.Nickname,
		PayeeID:     member.ID,
		PayeeCode:   member.PayeeCode,
	}
	app.Result(c, resp)
}

type ConfirmTransferReq struct {
	OrderNo   string          `json:"order_no"`
	Payee     int64           `json:"payee"`
	PayeeCode string          `json:"payee_code"`
	Currency  currency.Code   `json:"currency"`
	Amount    decimal.Decimal `json:"amount"`
	Desc      string          `json:"desc"`
	PIN       string          `json:"pin"`
}

func ConfirmTransfer(c *gin.Context) {
	req := new(ConfirmTransferReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if !radx.ValidTransferOrderNo(req.OrderNo) {
		app.InvalidParams(c, "invalid order number")
		return
	}
	args := &contracts.MemberTransferReq{
		PIN:   req.PIN,
		Payee: req.Payee,
		TransferReqCommon: contracts.TransferReqCommon{
			OrderNo:   req.OrderNo,
			From:      auth.MemberID(c),
			PayeeCode: req.PayeeCode,
			Reason:    ledger.Transfer,
			Currency:  currency.Coin,
			Amount:    req.Amount,
			Desc:      "平台用户转账",
		},
	}
	resp, err := trade_rpcx.MemberTransfer(c.Request.Context(), args)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
