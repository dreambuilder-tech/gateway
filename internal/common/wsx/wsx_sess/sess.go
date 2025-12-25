package wsx_sess

import "github.com/olahol/melody"

const key = "session_state"

type ChatInfo struct {
	PeerIds    map[int64]struct{}
	LastSendAt int64
}

type SessionState struct {
	Chats map[string]*ChatInfo // sessionID -> ChatInfo
}

func GetSessionState(s *melody.Session) *SessionState {
	v, ok := s.Get(key)
	if !ok {
		st := &SessionState{
			Chats: make(map[string]*ChatInfo),
		}
		s.Set(key, st)
		return st
	}
	st, _ := v.(*SessionState)
	if st.Chats == nil {
		st.Chats = make(map[string]*ChatInfo)
	}
	return st
}

func GetChatInfo(s *melody.Session, sessionID string) (*ChatInfo, bool) {
	state := GetSessionState(s)
	chat, ok := state.Chats[sessionID]
	return chat, ok
}

func SetChatInfo(s *melody.Session, sessionID string, info *ChatInfo) {
	state := GetSessionState(s)
	state.Chats[sessionID] = info
}
