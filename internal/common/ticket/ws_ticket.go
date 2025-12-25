package ticket

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
	"wallet/common-lib/env"

	"github.com/redis/go-redis/v9"
)

var (
	MalformedTicket = errors.New("malformed ticket")
	ExpiredTicket   = errors.New("expired ticket")
)

func websocketTicket(ticket string) string {
	return fmt.Sprintf("app:websocket:ticket:%s", ticket)
}

func randStr() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Set(ctx context.Context, rds redis.UniversalClient, sid string) (string, error) {
	str, err := randStr()
	if err != nil {
		return "", fmt.Errorf("create a new ticket error: %s", err.Error())
	}
	ttl := time.Minute
	if env.Prod() {
		ttl = time.Second * 5
	}
	ok, err := rds.SetNX(ctx, websocketTicket(str), sid, ttl).Result()
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("ticket collision")
	}
	return str, nil
}

func Use(ctx context.Context, rds redis.UniversalClient, ticket string) (string, error) {
	if !validTicket(ticket) {
		return "", MalformedTicket
	}
	r := rds.GetDel(ctx, websocketTicket(ticket))
	if err := r.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ExpiredTicket
		}
		return "", err
	}
	v := r.Val()
	if v == "" {
		return "", ExpiredTicket
	}
	return v, nil
}

func validTicket(t string) bool {
	if len(t) != 43 {
		return false
	}
	for _, ch := range t {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			continue
		}
		return false
	}
	return true
}
