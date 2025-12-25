package hubx

import (
	"context"
	"time"
	"wallet/common-lib/zapx"

	"github.com/gorilla/websocket"
	"github.com/olahol/melody"
	"go.uber.org/zap"
)

const (
	HeartbeatKey = "lastPong"

	heartbeatCheckInterval = 30 * time.Second
	heartbeatTimeout       = 300 * time.Second
)

func (h *Hub) ScanHeartbeat(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(heartbeatCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.scanThenCloseTimeoutSessions()
			}
		}
	}()
}

func (h *Hub) scanThenCloseTimeoutSessions() {
	now := time.Now()

	h.sessionStore.RangeAll(func(uid int64, s *melody.Session) bool {
		v, ok := s.Get(HeartbeatKey)
		if !ok {
			h.closeOnHeartbeatTimeout(uid, s, "no pong received yet")
			return true
		}
		last, ok := v.(time.Time)
		if !ok {
			h.closeOnHeartbeatTimeout(uid, s, "unsupported last pong type")
			return true
		}

		if now.Sub(last) > heartbeatTimeout {
			// trigger HandleDisconnect.
			h.closeOnHeartbeatTimeout(uid, s, "heartbeat timeout")
		}

		return true
	})
}

func (h *Hub) closeOnHeartbeatTimeout(memberID int64, s *melody.Session, reason string) {
	payload := websocket.FormatCloseMessage(4000, reason)
	if err := s.CloseWithMsg(payload); err != nil {
		zapx.ErrorCtx(context.Background(), "close session on heartbeat timeout error",
			zap.Int64("member", memberID),
			zap.String("reason", reason),
			zap.Error(err),
		)
		return
	}
	zapx.InfoCtx(context.Background(), "close session on heartbeat timeout",
		zap.Int64("member", memberID),
		zap.String("reason", reason),
	)
}
