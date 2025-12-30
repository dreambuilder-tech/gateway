package member

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/sms"

	"github.com/gin-gonic/gin"
)

type VerifyPhoneReq struct {
	AreaCode string `json:"area_code"`
	Phone    string `json:"phone"`
}

func VerifyPhone(c *gin.Context) {
	req := new(VerifyPhoneReq)
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
	resp, err := member_rpcx.VerifyPhone(c.Request.Context(), &contracts.VerifyPhoneReq{
		MemberID: auth.MemberID(c),
		AreaCode: req.AreaCode,
		Phone:    req.Phone,
	})
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

type PhoneBindReq struct {
	AreaCode string `json:"area_code"`
	Phone    string `json:"phone"`
	SMSCode  string `json:"sms_code"`
}

func PhoneBind(c *gin.Context) {
	req := new(PhoneBindReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.AreaCode == "" || req.Phone == "" || req.SMSCode == "" {
		app.InvalidParams(c, "empty parameters")
		return
	}
	resp, err := member_rpcx.PhoneBind(c.Request.Context(), &contracts.PhoneBindReq{
		AreaCode: req.AreaCode,
		Phone:    req.Phone,
		SMSCode:  req.SMSCode,
		MemberID: auth.MemberID(c),
	})
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
