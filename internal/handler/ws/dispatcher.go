package ws

import (
	"encoding/json"
	"fmt"
	"gateway/internal/common/wsx"
	"gateway/internal/common/wsx/hubx"
	"gateway/internal/common/wsx/ratex"

	// "gateway/internal/handler/ws/chat"

	"gateway/internal/handler/ws/system"
	wsproto "wallet/common-lib/wsx/proto"
	"wallet/common-lib/zapx"

	"github.com/olahol/melody"
	"go.uber.org/zap"
)

type Dispatcher struct {
	Hub       *hubx.Hub
	SystemSvr *system.Svr
	// ChatSvr   *chat.Svr
	Limiter *ratex.Limiter
}

func (d *Dispatcher) Dispatch(s *melody.Session, memberID int64, msg []byte) (*wsproto.PushMsg, error) {
	up := new(wsproto.UpMsg)
	if err := json.Unmarshal(msg, up); err != nil {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), fmt.Errorf("unmarshal up msg error: %w", err)
	}
	if !d.Limiter.Allow(memberID) {
		zapx.WarnCtx(s.Request.Context(), "ws message rate limited", zap.Int64("member", memberID), zap.String("from", s.Request.RemoteAddr))
		return wsx.SystemError(memberID, wsproto.CodeRateLimited), nil
	}
	switch up.Biz {
	case wsproto.BizSystem:
		return d.SystemSvr.Handle(s, memberID, up.Event, up.Payload)
	// case wsproto.BizChat:
	// 	return d.ChatSvr.Handle(s, memberID, up.Event, up.Payload)
	default:
		return nil, fmt.Errorf("unsupported biz: %s", up.Biz)
	}
}
