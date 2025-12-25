package router

import (
	"errors"
	"gateway/internal/app"
	"gateway/internal/common/auth"
	"gateway/internal/common/ticket"
	"gateway/internal/handler/ws"
	"wallet/common-lib/rdb"
	"wallet/common-lib/zapx"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"go.uber.org/zap"
)

func wsRouter(r *gin.RouterGroup) {
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		t := c.Query("ticket")
		if t == "" {
			app.Unauthorized(c, "empty ticket")
			return
		}
		sid, err := ticket.Use(c.Request.Context(), rdb.Client, t)
		if err != nil {
			if errors.Is(err, ticket.MalformedTicket) {
				app.InvalidParams(c, "wrong ticket")
				return
			}
			if errors.Is(err, ticket.ExpiredTicket) {
				app.RequestExpired(c, "expired ticket")
				return
			}
			zapx.ErrorCtx(c.Request.Context(), "use ticket error", zap.Error(err))
			app.InternalError(c, "use ticket error")
			return
		}
		if ok := auth.ParseSessionID(c, sid); !ok {
			return
		}
		memberID := auth.MemberID(c)
		keys := map[string]any{
			"memberID": memberID,
		}
		if err := m.HandleRequestWithKeys(c.Writer, c.Request, keys); err != nil {
			zapx.ErrorCtx(c.Request.Context(), "handle request error", zap.Int64("member", memberID), zap.Error(err))
		}
	})
	ws.Init(m)
}
