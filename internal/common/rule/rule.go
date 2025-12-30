package rule

import (
	"time"
	"wallet/common-lib/mw/biz_rule"
	"wallet/common-lib/utils/rate_limit"
)

var (
	Tourist = biz_rule.Rule{Rule: rate_limit.Rule{Window: time.Minute, MaxCount: 200}}
	Member  = biz_rule.Rule{Rule: rate_limit.Rule{Window: time.Minute, MaxCount: 300}}
)

var (
	SendSMS   = biz_rule.Rule{BizKey: "sms_send", Rule: rate_limit.Rule{Window: time.Minute, MaxCount: 3}}
	Login     = biz_rule.Rule{BizKey: "login", Rule: rate_limit.Rule{Window: time.Minute, MaxCount: 20}}
	Register  = biz_rule.Rule{BizKey: "register", Rule: rate_limit.Rule{Window: time.Minute, MaxCount: 20}}
	ChangePIN = biz_rule.Rule{BizKey: "change_pin", Rule: rate_limit.Rule{Window: time.Minute, MaxCount: 3}}
	Withdraw  = biz_rule.Rule{BizKey: "withdraw", Rule: rate_limit.Rule{Window: time.Minute, MaxCount: 3}}
)

var (
	ChatSendImage = biz_rule.Rule{BizKey: "chat_send_image", Rule: rate_limit.Rule{Window: time.Minute, MaxCount: 30, MaxSuccess: 10}}
)
