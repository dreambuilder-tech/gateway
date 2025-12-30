package middleware

import (
	"context"
	"fmt"
	"gateway/internal/common/auth"
	"gateway/internal/common/rule"
	"wallet/common-lib/app"
	"wallet/common-lib/mw/biz_rule"
	"wallet/common-lib/utils/rate_limit"
	"wallet/common-lib/zapx"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func IPLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := rateLimitSubject(c)
		allow, _ := checkLimit(c.Request.Context(), subject, rule.Tourist)
		if !allow {
			app.TooManyRequest(c)
			return
		}
		c.Next()
	}
}

func WithLimit(br biz_rule.Rule) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := rateLimitSubjectV2(c)
		allow, reserved := checkLimit(c.Request.Context(), subject, br)
		if !allow {
			app.TooManyRequest(c)
			return
		}
		c.Next()
		// if it has success limit.
		if reserved && br.Rule.MaxSuccess > 0 && !app.IsSuccess(c) {
			// rollback success count.
			rate_limit.Limiter.RollbackSuccess(c.Request.Context(), subject, br.BizKey)
		}
	}
}

func checkLimit(ctx context.Context, subject string, br biz_rule.Rule) (bool, bool) {
	if rate_limit.Limiter == nil {
		return true, false
	}
	ok, reserved, err := rate_limit.Limiter.Allow(ctx, subject, br.BizKey, br.Rule)
	if err != nil {
		zapx.ErrorCtx(ctx, "check rate limit error", zap.Error(err))
		return true, false
	}
	return ok, reserved
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
