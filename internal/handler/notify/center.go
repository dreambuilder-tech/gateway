package notify

import (
	"gateway/internal/common/auth"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/notify_mode"
	"wallet/common-lib/dto/req_dto"
	"wallet/common-lib/rpcx/im_rpcx"

	"github.com/gin-gonic/gin"
)

func CenterUnread(c *gin.Context) {
	resp, err := im_rpcx.CenterUnread(c.Request.Context(), auth.MemberID(c))
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

type CenterListReq struct {
	Mode notify_mode.Code `form:"mode" binding:"required"`
	req_dto.PageArgs
}

func CenterList(c *gin.Context) {
	req := new(CenterListReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	req.PageArgs.Init()
	resp, err := im_rpcx.CenterNotifications(c.Request.Context(), auth.MemberID(c), req.Mode, req.Page, req.Size)
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.ResultPage(c, resp.Data, resp.Total)
}
