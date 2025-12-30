package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
	"wallet/common-lib/app"
	"wallet/common-lib/rdb"
	"wallet/common-lib/zapx"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Member struct {
	ID       int64  `json:"id"`
	Account  string `json:"account"`
	DeviceID string `json:"device_id"`
	LoginIP  string `json:"login_ip"`
	LoginAt  int64  `json:"login_at"`
}

const (
	SessionHeader     = "X-Session-Id"
	ReqMemberID       = "memberID"
	ReqMemberAccount  = "memberAccount"
	SessionExpireTime = 7 * 24 * time.Hour
)

func NewSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GetSessionID(c *gin.Context) string {
	return c.GetHeader(SessionHeader)
}

func SetSessionID(c *gin.Context, sid string) {
	c.Header(SessionHeader, sid)
}

func MemberID(c *gin.Context) int64 {
	return c.GetInt64(ReqMemberID)
}

func MemberAccount(c *gin.Context) string {
	return c.GetString(ReqMemberAccount)
}

func ParseSessionID(c *gin.Context, sid string) bool {
	if sid == "" {
		app.Unauthorized(c, "empty session")
		return false
	}
	var (
		rds = rdb.Client
		ctx = c.Request.Context()
	)
	data, err := rds.Get(ctx, sid).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			app.Unauthorized(c, "session expired")
		} else {
			app.Unauthorized(c, "auth error")
		}
		zapx.ErrorCtx(ctx, "read session cache error", zap.Error(err))
		return false
	}
	member := new(Member)
	if err = json.Unmarshal(data, member); err != nil {
		zapx.ErrorCtx(ctx, "unmarshal session cache error", zap.Error(err))
		app.Unauthorized(c, "invalid session data")
		return false
	}
	_ = rds.Expire(ctx, sid, SessionExpireTime)
	c.Set(ReqMemberID, member.ID)
	c.Set(ReqMemberAccount, member.Account)
	return true
}
