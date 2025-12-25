package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/internal/common/wsx"
	"gateway/internal/common/wsx/hubx"
	"gateway/internal/common/wsx/wsx_sess"
	"time"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/im_rpcx"
	wsproto "wallet/common-lib/wsx/proto"
	"wallet/common-lib/wsx/wstopic"
	"wallet/common-lib/zapx"

	"github.com/olahol/melody"
	"go.uber.org/zap"
)

type Svr struct {
	Hub *hubx.Hub
}

func (c *Svr) Handle(s *melody.Session, memberID int64, event wsproto.Event, payload json.RawMessage) (*wsproto.PushMsg, error) {
	switch event {
	case wsproto.UpEventOrderChatJoin:
		return c.orderChatJoin(s, memberID, payload)
	case wsproto.UpEventOrderChatLeave:
		return c.orderChatLeave(s, memberID, payload)
	case wsproto.UpEventMsgSend:
		return c.messageSendAck(s, memberID, payload)
	case wsproto.UpEventMsgRead:
		return c.messageRead(s, memberID, payload)
	default:
		return nil, fmt.Errorf("unsupported event: %s", event)
	}
}

func (c *Svr) orderChatJoin(s *melody.Session, memberID int64, payload json.RawMessage) (*wsproto.PushMsg, error) {
	req := new(wsproto.JoinOrderChatReq)
	if err := json.Unmarshal(payload, req); err != nil {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), fmt.Errorf("unmarshal order chat req error: %w", err)
	}
	if req.SessionID == "" {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), errors.New("empty parameters")
	}
	ack := &wsproto.PushMsg{
		TargetK: wsproto.TargetUID,
		TargetV: memberID,
		Biz:     wsproto.BizChat,
		Event:   wsproto.DownEventOrderChatJoinAck,
		Code:    wsproto.CodeSuccess,
		Payload: &wsproto.JoinOrderChatAck{
			SessionID: req.SessionID,
		},
	}
	resp, err := im_rpcx.CanJoin(s.Request.Context(), &contracts.CanJoinReq{
		SessionID: req.SessionID,
		MemberID:  memberID,
	})
	if err != nil {
		ack.Code = wsproto.CodeInternalError
		return ack, fmt.Errorf("call rpcx error: %w", err)
	}
	if !resp.OK {
		ack.Code = wsproto.CodeNoPerms
		return ack, fmt.Errorf("member has no permission")
	}
	// store all peers to session state.
	sess := &wsx_sess.ChatInfo{
		PeerIds: make(map[int64]struct{}),
	}
	for _, v := range resp.PeerMembers {
		sess.PeerIds[v] = struct{}{}
	}
	wsx_sess.SetChatInfo(s, req.SessionID, sess)
	for k := range sess.PeerIds {
		c.Hub.Subscribe(s, memberID, wstopic.PeerOnlineTopic(k))
	}
	c.Hub.Subscribe(s, memberID, wstopic.ChatSessionTopic(req.SessionID))

	zapx.InfoCtx(s.Request.Context(), "member join order chat success", zap.Int64("uid", memberID), zap.String("session_id", req.SessionID))

	return ack, nil
}

func (c *Svr) orderChatLeave(s *melody.Session, memberID int64, payload json.RawMessage) (*wsproto.PushMsg, error) {
	req := new(wsproto.LeaveOrderChatReq)
	if err := json.Unmarshal(payload, req); err != nil {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), fmt.Errorf("unmarshal leave order chat msg fail: %w", err)
	}
	if req.SessionID == "" {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), errors.New("empty parameters")
	}
	c.Hub.Unsubscribe(s, memberID, wstopic.ChatSessionTopic(req.SessionID))

	sess := wsx_sess.GetSessionState(s)
	defer delete(sess.Chats, req.SessionID)
	if chatInfo, ok := sess.Chats[req.SessionID]; ok {
		for peer := range chatInfo.PeerIds {
			var flag bool
			for k, v := range sess.Chats {
				if k == req.SessionID {
					continue
				}
				if _, ok := v.PeerIds[peer]; ok {
					flag = true
					break
				}
			}
			if !flag {
				c.Hub.Unsubscribe(s, memberID, wstopic.PeerOnlineTopic(peer))
			}
		}
	}
	return nil, nil
}

func (c *Svr) messageSendAck(s *melody.Session, memberID int64, payload json.RawMessage) (*wsproto.PushMsg, error) {
	req := new(wsproto.SendMsgReq)
	if err := json.Unmarshal(payload, req); err != nil {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), fmt.Errorf("unmarshal send msg req error: %w", err)
	}
	if req.SessionID == "" {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), errors.New("empty parameters")
	}

	r := &wsproto.PushMsg{
		TargetK: wsproto.TargetUID,
		TargetV: memberID,
		Biz:     wsproto.BizChat,
		Event:   wsproto.DownEventMsgSendAck,
		Code:    wsproto.CodeSuccess,
	}
	ack := &wsproto.SendMsgAck{
		SessionID:   req.SessionID,
		ClientMsgID: req.ClientMsgID,
	}
	r.Payload = ack

	if code := sendLimit(s, req.SessionID); code != wsproto.CodeSuccess {
		r.Code = code
		return r, nil
	}

	if code := verifyContent(req.Content); code != wsproto.CodeSuccess {
		r.Code = code
		return r, fmt.Errorf("verify message content fail. code: %d", code)
	}

	ctx, cancel := context.WithTimeout(s.Request.Context(), 5*time.Second)
	defer cancel()
	args := &contracts.SendMsgReq{
		ClientMsgID: req.ClientMsgID,
		SessionID:   req.SessionID,
		SenderUID:   memberID,
		Content:     req.Content,
		SendAt:      time.Now().UnixMilli(),
	}
	resp, err := im_rpcx.SendMsg(ctx, args)
	if err != nil {
		r.Code = wsproto.CodeInternalError
		return r, fmt.Errorf("rpcx send message error: %w", err)
	}
	ack.MsgID = resp.MsgID
	return r, nil
}

func (c *Svr) messageRead(s *melody.Session, memberID int64, payload json.RawMessage) (*wsproto.PushMsg, error) {
	req := new(wsproto.ReadMsgReq)
	if err := json.Unmarshal(payload, req); err != nil {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), fmt.Errorf("unmarshal read msg req error: %w", err)
	}
	if req.SessionID == "" || req.LastMsgID <= 0 {
		return wsx.SystemError(memberID, wsproto.CodeInvalidParameter), errors.New("empty parameters")
	}
	args := &contracts.ReadMsgReq{
		SessionID: req.SessionID,
		ReaderUID: memberID,
		LastMsgID: req.LastMsgID,
	}
	_, err := im_rpcx.ReadMsg(s.Request.Context(), args)
	if err != nil {
		return wsx.SystemError(memberID, wsproto.CodeInternalError), fmt.Errorf("rpcx read message error: %w", err)
	}
	return nil, nil
}
