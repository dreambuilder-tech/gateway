package chat

import (
	"gateway/internal/app"
	"gateway/internal/common/auth"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/im_rpcx"

	"github.com/gin-gonic/gin"
)

type HistoryReq struct {
	SessionID string `json:"session_id"`
	LastMsgID int64  `json:"last_msg_id"`
	Size      int    `json:"size"`
}

func History(c *gin.Context) {
	req := new(HistoryReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.SessionID == "" {
		app.InvalidParams(c, "empty session id")
		return
	}
	resp, err := im_rpcx.ChatHistory(c.Request.Context(), &contracts.ChatHistoryReq{
		SessionID: req.SessionID,
		MemberID:  auth.MemberID(c),
		LastMsgID: req.LastMsgID,
		Size:      req.Size,
	})
	if err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	app.SuccessData(c, resp)
}
