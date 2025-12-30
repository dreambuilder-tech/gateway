package member

import (
	"gateway/internal/common/auth"
	"github.com/gin-gonic/gin"
	"wallet/common-lib/app"
	"wallet/common-lib/rpcx/member_rpcx"
)

func DeviceLoginRecordList(c *gin.Context) {
	resp, err := member_rpcx.DeviceLoginRecordList(c.Request.Context(), auth.MemberID(c))
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}
