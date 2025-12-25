package chat

import (
	"gateway/internal/app"
	"gateway/internal/common/auth"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/im_rpcx"

	"github.com/gin-gonic/gin"
)

type InitReq struct {
	SessionID string `json:"session_id"`
	Size      int    `json:"size"`
}

func Init(c *gin.Context) {
	req := new(InitReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.SessionID == "" {
		app.InvalidParams(c, "empty session id")
		return
	}
	resp, err := im_rpcx.InitChat(c.Request.Context(), &contracts.InitChatReq{
		SessionID: req.SessionID,
		MemberID:  auth.MemberID(c),
		Size:      req.Size,
	})
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	app.SuccessData(c, resp)
}
