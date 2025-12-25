package hubx

import "sync/atomic"

type Stats struct {
	connect    atomic.Int64
	disconnect atomic.Int64
	msgIn      atomic.Int64
}

func (s *Stats) IncConnect() {
	s.connect.Add(1)
}

func (s *Stats) IncDisconnect() {
	s.disconnect.Add(1)
}

func (s *Stats) IncMsgIn() {
	s.msgIn.Add(1)
}

func (s *Stats) Snapshot() (c, d, m int64) {
	return s.connect.Load(), s.disconnect.Load(), s.msgIn.Load()
}
