package rule

import (
	"time"
	"wallet/common-lib/utils/rate_limit"
)

type BizRule struct {
	BizKey string
	Rule   rate_limit.Rule
}

var (
	Tourist = BizRule{"", rate_limit.Rule{Window: time.Minute, MaxCount: 200}}
	Member  = BizRule{"", rate_limit.Rule{Window: time.Minute, MaxCount: 300}}
)

var (
	Login    = BizRule{"login", rate_limit.Rule{Window: time.Minute, MaxCount: 20}}
	Register = BizRule{"register", rate_limit.Rule{Window: time.Minute, MaxCount: 20}}
	SendSMS  = BizRule{"sms_send", rate_limit.Rule{Window: time.Minute, MaxCount: 3}}
	Withdraw = BizRule{"withdraw", rate_limit.Rule{Window: time.Minute, MaxCount: 3}}
)
