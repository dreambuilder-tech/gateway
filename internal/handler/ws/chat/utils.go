package chat

import (
	"gateway/internal/common/wsx/wsx_sess"
	"time"
	"unicode/utf8"
	"wallet/common-lib/consts"
	"wallet/common-lib/utils/stringx"
	wsproto "wallet/common-lib/wsx/proto"

	"github.com/olahol/melody"
)

func sendLimit(s *melody.Session, sessionID string) wsproto.Code {
	now := time.Now().UnixMilli()
	chat, ok := wsx_sess.GetChatInfo(s, sessionID)
	if !ok {
		return wsproto.CodeNoPerms
	}
	if (now - chat.LastSendAt) < consts.ChatSendIntervalMS {
		return wsproto.CodeRateLimited
	}
	chat.LastSendAt = now
	return wsproto.CodeSuccess
}

func verifyContent(content *wsproto.Content) wsproto.Code {
	if content == nil {
		return wsproto.CodeContentEmpty
	}
	var (
		text   = content.Text
		images = content.Images
	)
	if text == "" && len(images) == 0 {
		return wsproto.CodeContentEmpty
	}
	if text != "" {
		if !stringx.IsSafe(text) {
			return wsproto.CodeTextMalformed
		}
		if utf8.RuneCountInString(text) > consts.ChatTextLengthLimit {
			return wsproto.CodeTextTooLong
		}
	}
	if len(images) > consts.ChatImagesLimit {
		return wsproto.CodeImageCountLimit
	}
	return wsproto.CodeSuccess
}
