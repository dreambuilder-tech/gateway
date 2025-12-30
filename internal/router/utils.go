package router

import (
	"gateway/internal/middleware"
	"wallet/common-lib/mw/biz_rule"

	"github.com/gin-gonic/gin"
)

func withLimit(br biz_rule.Rule) gin.HandlerFunc {
	return middleware.WithLimit(br)
}
