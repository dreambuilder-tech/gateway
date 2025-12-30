package member

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/sms"

	"github.com/gin-gonic/gin"
)

type CommitPinReq struct {
	OldPIN string `json:"old_pin"`
	NewPIN string `json:"new_pin"`
}

func CommitPin(c *gin.Context) {
	req := new(CommitPinReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	args := &contracts.CommitPinReq{
		MemberID: auth.MemberID(c),
		OldPIN:   req.OldPIN,
		NewPIN:   req.NewPIN,
	}
	resp, err := member_rpcx.CommitPIN(c.Request.Context(), args)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

type ResetPinReq struct {
	AreaCode string `json:"area_code"`
	Phone    string `json:"phone"`
	SMSCode  string `json:"sms_code"`
}

func ResetPin(c *gin.Context) {
	req := new(ResetPinReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if !sms.ValidArea(req.AreaCode) {
		app.InvalidParams(c, "invalid area code")
		return
	}
	if !sms.ValidPhone(req.Phone) {
		app.InvalidParams(c, "invalid phone number")
		return
	}
	args := &contracts.ResetPinReq{
		MemberID: auth.MemberID(c),
		AreaCode: req.AreaCode,
		Phone:    req.Phone,
		SMSCode:  req.SMSCode,
	}
	resp, err := member_rpcx.ResetPIN(c.Request.Context(), args)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

type ChangePinReq struct {
	Hash   string `json:"hash"`
	NewPIN string `json:"new_pin"`
}

func ChangePIN(c *gin.Context) {
	req := new(ChangePinReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	args := &contracts.ChangePinReq{
		MemberID: auth.MemberID(c),
		Hash:     req.Hash,
		NewPIN:   req.NewPIN,
	}
	resp, err := member_rpcx.ChangePIN(c.Request.Context(), args)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
