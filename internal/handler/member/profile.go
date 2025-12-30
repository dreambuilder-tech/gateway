package member

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/urlx"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	resp, err := member_rpcx.MemberProfile(c.Request.Context(), auth.MemberID(c))
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

type ChangeAvatarReq struct {
	Avatar string `json:"avatar"`
}

func ChangeAvatar(c *gin.Context) {
	req := new(ChangeAvatarReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.Avatar == "" {
		app.InvalidParams(c, "empty avatar")
		return
	}
	// 目前不允许自定义上传，支持自定义上传后也需要后台审核
	if urlx.IsValid(req.Avatar) || len(req.Avatar) > 50 {
		app.InvalidParams(c, "invalid avatar")
		return
	}
	resp, err := member_rpcx.MemberChangeAvatar(c.Request.Context(), auth.MemberID(c), req.Avatar)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
