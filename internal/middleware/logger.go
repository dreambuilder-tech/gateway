package middleware

import (
	"time"
	"wallet/common-lib/utils/timex"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const slowLogTime = 1800 * time.Millisecond

func LoggerWithZap(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency, slow := timex.GetLatency(start, slowLogTime)
		fs := []zap.Field{
			zap.String("ip", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("protocol", c.Request.Proto),
			zap.String("content-type", c.GetHeader("Content-Type")),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("query", query),
			zap.String("latency", latency.String()),
			zap.Int("slow", slow),
		}
		if traceID := c.GetHeader("X-Trace-Id"); traceID != "" {
			zap.String("trace-id", traceID)
		}
		if len(c.Errors) > 0 {
			fs = append(fs, zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
			l.Error(path, fs...)
		} else {
			l.Info(path, fs...)
		}
	}
}
