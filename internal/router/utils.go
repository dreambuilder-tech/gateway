package router

import (
	"gateway/internal/common/rule"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func withLimit(br rule.BizRule) gin.HandlerFunc {
	return middleware.WithLimit(br)
}
