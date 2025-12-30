package agent

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"

	"github.com/gin-gonic/gin"
)

func Preflight(c *gin.Context) {
	resp, err := member_rpcx.AgentApplyPreflight(c.Request.Context(), auth.MemberID(c))
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

type ApplyReq struct {
	PIN string `json:"pin"`
}

func Apply(c *gin.Context) {
	req := new(ApplyReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	resp, err := member_rpcx.AgentApply(c.Request.Context(), &contracts.AgentApplyReq{
		MemberID: auth.MemberID(c),
		PIN:      req.PIN,
	})
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

type RefundReq struct {
	PIN string `json:"pin"`
}

func Refund(c *gin.Context) {
	req := new(RefundReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	resp, err := member_rpcx.AgentRefund(c.Request.Context(), &contracts.AgentRefundReq{
		MemberID: auth.MemberID(c),
		PIN:      req.PIN,
	})
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
