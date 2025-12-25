package middleware

import (
	"gateway/internal/common/auth"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := auth.GetSessionID(c)
		if ok := auth.ParseSessionID(c, sid); !ok {
			return
		}
		c.Next()
	}
}
