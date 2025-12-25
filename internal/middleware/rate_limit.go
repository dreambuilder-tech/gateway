package middleware

import (
	"fmt"
	"gateway/internal/app"
	"gateway/internal/common/auth"
	"gateway/internal/common/rule"
	"wallet/common-lib/utils/rate_limit"
	"wallet/common-lib/zapx"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func GlobalLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := rateLimitSubject(c)
		if !checkLimit(c, subject, rule.Tourist) {
			app.TooManyRequest(c)
			return
		}
		c.Next()
	}
}

func WithLimit(br rule.BizRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := rateLimitSubjectV2(c)
		if !checkLimit(c, subject, br) {
			app.TooManyRequest(c)
			return
		}
		c.Next()
	}
}

func checkLimit(c *gin.Context, subject string, br rule.BizRule) bool {
	if rate_limit.Limiter == nil {
		return true
	}
	ok, err := rate_limit.Limiter.Allow(c.Request.Context(), subject, br.BizKey, br.Rule)
	if err != nil {
		zapx.ErrorCtx(c.Request.Context(), "check rate limit error", zap.Error(err))
		return true
	}
	return ok
}

func rateLimitSubject(c *gin.Context) string {
	return fmt.Sprintf("ip:%v", c.ClientIP())
}

func rateLimitSubjectV2(c *gin.Context) string {
	if uid, ok := c.Get(auth.ReqMemberID); ok {
		return fmt.Sprintf("uid:%v", cast.ToString(uid))
	} else {
		return fmt.Sprintf("ip:%v", c.ClientIP())
	}
}
