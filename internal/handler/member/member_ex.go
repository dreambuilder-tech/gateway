package member

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/stringx"

	"github.com/gin-gonic/gin"
)

// only used to bindJson from FE
type ChangeNicknameReq struct {
	NewNickname string `json:"new_nickname"`
}

func ChangeNickname(c *gin.Context) {
	req := new(ChangeNicknameReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}

	// check
	if req.NewNickname == "" {
		app.InvalidParams(c, "empty parameters")
		return
	}

	// 计算中文长度 如"123威廉"，返回5
	countIncCHN := stringx.GetRuneCount(req.NewNickname)
	if countIncCHN < 2 || countIncCHN > 8 || !stringx.IsChineseEnglishNumber(req.NewNickname) {
		app.InvalidParams(c, "invalid parameters")
		return
	}

	resp, err := member_rpcx.ChangeNickname(c.Request.Context(), &contracts.ChangeNickNameReq{
		NewNickname: req.NewNickname,
		MemberID:    auth.MemberID(c),
	})
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
