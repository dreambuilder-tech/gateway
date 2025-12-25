package hubx

import (
	"sync"

	"github.com/olahol/melody"
)

type SessionStore struct {
	mu sync.RWMutex

	// uid -> sessions
	userSubs map[int64]map[*melody.Session]struct{}

	// topic -> sessions
	topicSubs map[string]map[*melody.Session]struct{}

	// session -> uid
	sessionUID map[*melody.Session]int64

	// session -> topics
	sessionTopics map[*melody.Session]map[string]struct{}
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		userSubs:      make(map[int64]map[*melody.Session]struct{}),
		topicSubs:     make(map[string]map[*melody.Session]struct{}),
		sessionUID:    make(map[*melody.Session]int64),
		sessionTopics: make(map[*melody.Session]map[string]struct{}),
	}
}

func (st *SessionStore) BindUID(s *melody.Session, uid int64) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.sessionUID[s] = uid

	if st.userSubs[uid] == nil {
		st.userSubs[uid] = make(map[*melody.Session]struct{})
	}
	st.userSubs[uid][s] = struct{}{}
}

func (st *SessionStore) Subscribe(s *melody.Session, topic string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if st.topicSubs[topic] == nil {
		st.topicSubs[topic] = make(map[*melody.Session]struct{})
	}
	st.topicSubs[topic][s] = struct{}{}

	if st.sessionTopics[s] == nil {
		st.sessionTopics[s] = make(map[string]struct{})
	}
	st.sessionTopics[s][topic] = struct{}{}
}

func (st *SessionStore) Unsubscribe(s *melody.Session, topic string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if m, ok := st.topicSubs[topic]; ok {
		delete(m, s)
		if len(m) == 0 {
			delete(st.topicSubs, topic)
		}
	}
	if ts, ok := st.sessionTopics[s]; ok {
		delete(ts, topic)
		if len(ts) == 0 {
			delete(st.sessionTopics, s)
		}
	}
}

func (st *SessionStore) RemoveSession(s *melody.Session) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if uid, ok := st.sessionUID[s]; ok {
		if ss, ok2 := st.userSubs[uid]; ok2 {
			delete(ss, s)
			if len(ss) == 0 {
				delete(st.userSubs, uid)
			}
		}
		delete(st.sessionUID, s)
	}

	if topics, ok := st.sessionTopics[s]; ok {
		for topic := range topics {
			if subs, ok2 := st.topicSubs[topic]; ok2 {
				delete(subs, s)
				if len(subs) == 0 {
					delete(st.topicSubs, topic)
				}
			}
		}
		delete(st.sessionTopics, s)
	}
}

func (st *SessionStore) GetByUID(uid int64) []*melody.Session {
	st.mu.RLock()
	defer st.mu.RUnlock()

	m := st.userSubs[uid]
	if len(m) == 0 {
		return nil
	}
	res := make([]*melody.Session, 0, len(m))
	for s := range m {
		res = append(res, s)
	}
	return res
}

func (st *SessionStore) GetByTopic(topic string) []*melody.Session {
	st.mu.RLock()
	defer st.mu.RUnlock()

	m := st.topicSubs[topic]
	if len(m) == 0 {
		return nil
	}
	res := make([]*melody.Session, 0, len(m))
	for s := range m {
		res = append(res, s)
	}
	return res
}

func (st *SessionStore) RangeAll(fn func(uid int64, sess *melody.Session) bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	for sess, uid := range st.sessionUID {
		if !fn(uid, sess) {
			return
		}
	}
}

func (st *SessionStore) sessionCount() int {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return len(st.sessionUID)
}

func (st *SessionStore) userSessionCount(uid int64) int {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return len(st.userSubs[uid])
}
