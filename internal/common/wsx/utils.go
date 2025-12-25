package wsx

import wsproto "wallet/common-lib/wsx/proto"

func SystemError(memberID int64, code wsproto.Code) *wsproto.PushMsg {
	return &wsproto.PushMsg{
		TargetK: wsproto.TargetUID,
		TargetV: memberID,
		Biz:     wsproto.BizSystem,
		Event:   wsproto.DownEventSystemError,
		Code:    code,
	}
}
