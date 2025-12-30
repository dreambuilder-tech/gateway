package member

import (
	"gateway/internal/common/auth"
	"github.com/gin-gonic/gin"
	"wallet/common-lib/app"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/stringx"
)

type VerifyRealNameReq struct {
	RealName string `json:"real_name"`
	IDNumber string `json:"id_number"`
}

func VerifyRealName(c *gin.Context) {
	req := new(VerifyRealNameReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.RealName == "" || req.IDNumber == "" {
		app.InvalidParams(c, "empty parameters")
		return
	}
	if !stringx.IsChinese(req.RealName, 20) || !stringx.IsIDCard(req.IDNumber) {
		app.InvalidParams(c, "invalid parameters")
		return
	}

	resp, err := member_rpcx.VerifyRealName(c.Request.Context(), &contracts.VerifyRealNameReq{
		RealName: req.RealName,
		IDNumber: req.IDNumber,
		MemberID: auth.MemberID(c),
	})
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
