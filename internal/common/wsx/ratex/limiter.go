package ratex

import (
	"sync"

	"golang.org/x/time/rate"
)

type Limiter struct {
	mu     sync.Mutex
	bucket map[int64]*rate.Limiter
	rate   rate.Limit
	burst  int
}

func NewLimiter(r rate.Limit, burst int) *Limiter {
	return &Limiter{
		bucket: make(map[int64]*rate.Limiter),
		rate:   r,
		burst:  burst,
	}
}

func (l *Limiter) Allow(memberID int64) bool {
	if memberID <= 0 {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	lim, ok := l.bucket[memberID]
	if !ok {
		lim = rate.NewLimiter(l.rate, l.burst)
		l.bucket[memberID] = lim
	}
	return lim.Allow()
}
