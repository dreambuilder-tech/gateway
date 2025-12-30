package chat

import (
	"context"
	"gateway/internal/common/wsx/wsx_sess"
	"time"
	"unicode/utf8"
	"wallet/common-lib/consts"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/stringx"
	wsproto "wallet/common-lib/wsx/proto"

	"github.com/olahol/melody"
	"go.uber.org/zap"
)

// sendLimit 用户(s)在当前房间(sessionID)的聊天频率
func sendLimit(s *melody.Session, sessionID string) wsproto.Code {
	now := time.Now().UnixMilli()
	chat, ok := wsx_sess.GetChatInfo(s, sessionID)
	if !ok {
		return wsproto.CodePermissionDenied
	}
	if (now - chat.LastSentAt) < consts.ChatSendIntervalMS {
		return wsproto.CodeRateLimited
	}
	chat.LastSentAt = now
	return wsproto.CodeSuccess
}

func verifyBlocks(memberID int64, blocks []*wsproto.Block) wsproto.Code {
	if blocks == nil || len(blocks) == 0 {
		return wsproto.CodeContentEmpty
	}
	if len(blocks) > consts.ChatBlocksLimit {
		return wsproto.CodeBlocksLimit
	}
	empty := true
	for _, v := range blocks {
		if v == nil {
			continue
		}
		switch v.Type {
		case wsproto.Text:
			text, ok := v.Data.(string)
			if !ok {
				return wsproto.CodeTextMalformed
			}
			if !stringx.IsSafe(text) {
				return wsproto.CodeTextMalformed
			}
			if utf8.RuneCountInString(text) > consts.ChatTextLengthLimit {
				return wsproto.CodeTextTooLong
			}
			empty = false
		case wsproto.Image:
			fileID, ok := v.Data.(int64)
			if !ok || fileID <= 0 {
				return wsproto.CodeImageFormatErr
			}
			// todo 目前场景只有单图上传
			resp, err := member_rpcx.FileCheckUse(context.Background(), &contracts.FileCheckUseReq{
				MemberID: memberID,
				FileID:   fileID,
			})
			if err != nil {
				zap.L().Error("image check err", zap.Error(err))
				continue
			}
			if resp.URL != "" {
				v.Data = resp.URL
				empty = false
			}
		}
	}
	if empty {
		return wsproto.CodeContentEmpty
	}
	return wsproto.CodeSuccess
}
