package ws

import (
	"context"
	"gateway/internal/common/wsx/hubx"
	"gateway/internal/common/wsx/ratex"
	"gateway/internal/ws/chat"
	"gateway/internal/ws/system"
	"time"
	"wallet/common-lib/natsx/event"
	"wallet/common-lib/rdb"
	"wallet/common-lib/rdb/domain/presence"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/zapx"

	"github.com/gorilla/websocket"
	"github.com/olahol/melody"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func Init(m *melody.Melody) {
	hub := hubx.New(zap.L().With(zap.String("unit", "ws-hub")))
	d := &Dispatcher{
		Hub: hub,
		ChatSvr: &chat.Svr{
			Hub: hub,
		},
		SystemSvr: &system.Svr{
			Hub: hub,
		},
		Limiter: ratex.NewLimiter(200, 400),
	}

	d.Hub.ScanHeartbeat(context.Background())
	d.Hub.TriggerStatsLogger(context.Background())

	d.Hub.SubscribeNatsTopic(event.ToChat)
	d.Hub.SubscribeNatsTopic(event.ToNotify)
	
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		memberID, ok := before(s)
		if !ok {
			return
		}
		d.Hub.Stats().IncMsgIn()
		zapx.DebugCtx(s.Request.Context(), "received message",
			zap.Int64("member", memberID),
			zap.String("from", s.Request.RemoteAddr),
			zap.String("message", string(msg)),
		)
		r, err := d.Dispatch(s, memberID, msg)
		if err != nil {
			zapx.ErrorCtx(s.Request.Context(), "dispatch failed", zap.Int64("member", memberID), zap.Error(err))
		}
		d.Hub.Push(r, s)
	})

	m.HandleConnect(func(s *melody.Session) {
		memberID, ok := before(s)
		if !ok {
			return
		}
		d.Hub.OnConnect(s, memberID)
		zapx.InfoCtx(s.Request.Context(), "ws connect", zap.Int64("member", memberID), zap.String("from", s.Request.RemoteAddr))
	})

	m.HandleDisconnect(func(s *melody.Session) {
		d.Hub.OnDisconnect(s)
		str, _ := s.Get("memberID")
		memberID := cast.ToInt64(str)
		if memberID > 0 {
			_ = member_rpcx.MemberWsClose(s.Request.Context(), memberID) // notify member-service -> fanout -> all gateway -> all sessions
		}
		zapx.InfoCtx(s.Request.Context(), "ws disconnect", zap.Any("member", memberID), zap.String("from", s.Request.RemoteAddr))
	})

	m.HandlePong(func(s *melody.Session) {
		before(s)
	})
}

func before(s *melody.Session) (int64, bool) {
	str, _ := s.Get("memberID")
	memberID := cast.ToInt64(str)
	if memberID <= 0 {
		_ = s.CloseWithMsg(websocket.FormatCloseMessage(4001, "unauthorized"))
		zapx.WarnCtx(s.Request.Context(), "user not login", zap.String("from", s.Request.RemoteAddr))
		return 0, false
	}
	s.Set(hubx.HeartbeatKey, time.Now())
	if err := presence.UpdateLastSeen(s.Request.Context(), rdb.Client, memberID); err != nil {
		zapx.ErrorCtx(s.Request.Context(), "presence update error", zap.Int64("member", memberID), zap.Error(err))
	}
	return memberID, true
}
