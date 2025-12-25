package hubx

import (
	"context"
	"encoding/json"
	"time"
	wsproto "wallet/common-lib/wsx/proto"
	"wallet/common-lib/zapx"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/olahol/melody"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const (
	statsLogInterval  = 30 * time.Second
	connectsLimit     = 100000
	userConnectsLimit = 15
)

type Hub struct {
	sessionStore *SessionStore
	logger       *zap.Logger
	stats        *Stats
}

func New(logger *zap.Logger) *Hub {
	return &Hub{
		sessionStore: NewSessionStore(),
		logger:       logger,
		stats:        &Stats{},
	}
}

func (h *Hub) Stats() *Stats {
	return h.stats
}

func (h *Hub) OnConnect(s *melody.Session, memberID int64) {
	if d := h.sessionStore.sessionCount(); d >= connectsLimit {
		zapx.WarnCtx(s.Request.Context(), "ws reject: server busy",
			zap.Int64("member", memberID), zap.String("from", s.Request.RemoteAddr),
			zap.Int("current", d), zap.Int("max", connectsLimit),
		)
		_ = s.CloseWithMsg(websocket.FormatCloseMessage(4002, "server busy"))
		return
	}
	if d := h.sessionStore.userSessionCount(memberID); d >= userConnectsLimit {
		zapx.WarnCtx(s.Request.Context(), "ws reject: user have too many connections",
			zap.Int64("member", memberID), zap.String("from", s.Request.RemoteAddr),
			zap.Int("current", d), zap.Int("max", userConnectsLimit),
		)
		_ = s.CloseWithMsg(websocket.FormatCloseMessage(4003, "too many connections"))
		return
	}
	h.stats.IncConnect()
	h.sessionStore.BindUID(s, memberID)
}

func (h *Hub) OnDisconnect(s *melody.Session) {
	h.stats.IncDisconnect()
	h.sessionStore.RemoveSession(s)
}

func (h *Hub) Subscribe(s *melody.Session, memberID int64, topics ...string) {
	for _, v := range topics {
		zapx.DebugCtx(s.Request.Context(), "subscribe topic", zap.Int64("member", memberID), zap.String("topic", v))
		h.sessionStore.Subscribe(s, v)
	}
}

func (h *Hub) Unsubscribe(s *melody.Session, memberID int64, topics ...string) {
	for _, v := range topics {
		zapx.DebugCtx(s.Request.Context(), "unsubscribe topic", zap.Int64("member", memberID), zap.String("topic", v))
		h.sessionStore.Unsubscribe(s, v)
	}
}

func (h *Hub) HandleNats(m *nats.Msg) {
	p := new(wsproto.PushMsg)
	if err := json.Unmarshal(m.Data, p); err != nil {
		h.logger.Error("HandleNats unmarshal error", zap.Error(err), zap.String("subject", m.Subject))
		return
	}
	if p.TargetK == wsproto.TargetSession {
		h.logger.Error("unsupported target from NATS")
		return
	}
	h.Push(p, nil)
}

func (h *Hub) Push(m *wsproto.PushMsg, s *melody.Session) {
	if m == nil {
		return
	}
	down := &wsproto.DownMsg{
		Biz:      m.Biz,
		Event:    m.Event,
		Code:     m.Code,
		ServerTs: time.Now().UnixMilli(),
		Payload:  m.Payload,
	}

	body, err := json.Marshal(down)
	if err != nil {
		h.logger.Error("marshal down data error", zap.Error(err))
		return
	}

	switch m.TargetK {
	case wsproto.TargetSession:
		if s != nil {
			if err = s.Write(body); err != nil {
				h.logger.Error("session write error", zap.Error(err))
			}
		}
	case wsproto.TargetUID:
		memberID := cast.ToInt64(m.TargetV)
		h.pushToUID(memberID, body)
	case wsproto.TargetTopic:
		topic := cast.ToString(m.TargetV)
		h.pushToTopic(topic, body)
	default:
		h.logger.Error("unsupported target", zap.Any("target", m.TargetK))
	}
}

func (h *Hub) TriggerStatsLogger(ctx context.Context) {
	ticker := time.NewTicker(statsLogInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c, d, m := h.stats.Snapshot()
				h.logger.Info("ws stats", zap.Int64("connect", c), zap.Int64("disconnect", d), zap.Int64("online", c-d), zap.Int64("msg_in", m))
			}
		}
	}()
}

func (h *Hub) pushToUID(memberID int64, msg []byte) {
	if memberID <= 0 {
		h.logger.Error("empty memberID")
		return
	}
	sessions := h.sessionStore.GetByUID(memberID)
	for _, s := range sessions {
		if err := s.Write(msg); err != nil {
			h.logger.Error("pushToUID write error", zap.Error(err))
			continue
		}
	}
}

func (h *Hub) pushToTopic(topic string, msg []byte) {
	if topic == "" {
		h.logger.Error("empty topic")
		return
	}
	sessions := h.sessionStore.GetByTopic(topic)
	for _, s := range sessions {
		if err := s.Write(msg); err != nil {
			h.logger.Error("pushToTopic write error", zap.Error(err))
			continue
		}
	}
}
