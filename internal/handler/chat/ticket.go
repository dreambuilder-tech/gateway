package chat

import (
	"gateway/internal/common/auth"
	"gateway/internal/common/ticket"
	"wallet/common-lib/app"
	"wallet/common-lib/rdb"
	"wallet/common-lib/zapx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TicketResp struct {
	Ticket string `json:"ticket"`
}

func Ticket(c *gin.Context) {
	sid := auth.GetSessionID(c)
	if sid == "" {
		app.Unauthorized(c, "empty session")
		return
	}
	ctx := c.Request.Context()
	t, err := ticket.Set(ctx, rdb.Client, sid)
	if err != nil {
		zapx.ErrorCtx(ctx, "set ticket error", zap.Error(err))
		app.InternalError(c, "set ticket error. try again")
		return
	}
	app.Result(c, &TicketResp{Ticket: t})
}
