package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func Response() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		if span != nil && span.SpanContext().IsValid() {
			xTraceId := span.SpanContext().TraceID().String()
			c.Header("X-Trace-Id", xTraceId)
		}
		c.Next()
	}
}
