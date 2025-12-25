package system

import (
	"encoding/json"
	"fmt"
	"gateway/internal/common/wsx/hubx"
	"time"
	wsproto "wallet/common-lib/wsx/proto"

	"github.com/olahol/melody"
)

type Svr struct {
	Hub *hubx.Hub
}

func (d *Svr) Handle(s *melody.Session, _ int64, event wsproto.Event, payload json.RawMessage) (*wsproto.PushMsg, error) {
	switch event {
	case wsproto.UpEventHeartbeat:
		return d.Heartbeat(s, payload)
	default:
		return nil, fmt.Errorf("unsupported event: %s", event)
	}
}

func (d *Svr) Heartbeat(s *melody.Session, payload json.RawMessage) (*wsproto.PushMsg, error) {
	req := new(wsproto.HeartbeatReq)
	if err := json.Unmarshal(payload, req); err != nil {
		return nil, err
	}
	s.Set(hubx.HeartbeatKey, time.Now())
	return &wsproto.PushMsg{
		TargetK: wsproto.TargetSession,
		Biz:     wsproto.BizSystem,
		Event:   wsproto.DownEventHeartbeatAck,
		Code:    wsproto.CodeSuccess,
		Payload: &wsproto.HeartbeatAck{
			ClientTs: req.ClientTs,
		},
	}, nil
}
