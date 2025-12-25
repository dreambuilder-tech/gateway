package member

import (
	"encoding/json"
	"gateway/internal/app"
	"gateway/internal/common/auth"
	"time"
	auth_cst "wallet/common-lib/consts/auth"
	"wallet/common-lib/rdb"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/sms"

	"github.com/gin-gonic/gin"
)

type SendSmsReq struct {
	Purpose auth_cst.SmsPurpose `json:"purpose"`
	Area    string              `json:"area"`
	Number  string              `json:"number"`
}

func SendSms(c *gin.Context) {
	req := new(SendSmsReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if !sms.ValidArea(req.Area) {
		app.InvalidParams(c, "invalid area code")
		return
	}
	if !sms.ValidPhone(req.Number) {
		app.InvalidParams(c, "invalid phone number")
		return
	}
	resp, err := member_rpcx.SendSms(c.Request.Context(), req.Purpose, req.Area, req.Number)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.SuccessData(c, resp)
}

func Register(c *gin.Context) {
	req := new(contracts.RegisterReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	resp, err := member_rpcx.Register(c.Request.Context(), req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.SuccessData(c, resp)
}

func Login(c *gin.Context) {
	req := new(contracts.MemberLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	resp, err := member_rpcx.Login(c.Request.Context(), req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	if resp.Fail() {
		app.SuccessData(c, resp)
		return
	}
	if err = storeSession(c, resp); err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.SuccessData(c, resp)
}

func ResetPwd(c *gin.Context) {
	req := new(contracts.ResetPwdReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.AreaCode == "" || req.Phone == "" || req.SMSCode == "" {
		app.InvalidParams(c, "empty parameters")
		return
	}
	resp, err := member_rpcx.ResetPwd(c.Request.Context(), req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.SuccessData(c, resp)
}

func ChangePwd(c *gin.Context) {
	req := new(contracts.ChangePwdReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.Hash == "" || req.AreaCode == "" || req.Phone == "" || req.Password == "" {
		app.InvalidParams(c, "empty parameters")
		return
	}
	resp, err := member_rpcx.ChangePwd(c.Request.Context(), req)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.SuccessData(c, resp)
}

func storeSession(c *gin.Context, m *contracts.MemberLoginResp) error {
	sid, err := auth.NewSessionID()
	if err != nil {
		return err
	}
	auth.SetSessionID(c, sid)
	user := &auth.Member{
		ID:      m.ID,
		LoginIP: c.ClientIP(),
		LoginAt: time.Now().Unix(),
	}
	bytes, _ := json.Marshal(user)
	return rdb.Client.Set(c, sid, bytes, auth.SessionExpireTime).Err()
}
