package notify

import (
	"context"
	"gateway/internal/common/auth"
	"strconv"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/ack_kind"
	"wallet/common-lib/rdb"
	"wallet/common-lib/rdb/domain/unread"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/im_rpcx"
	"wallet/common-lib/wsx/proto"
	"wallet/common-lib/zapx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Sync(c *gin.Context) {
	resp, err := im_rpcx.SyncNotify(c.Request.Context(), auth.MemberID(c))
	if err != nil {
		app.InternalError(c, err.Error())
		return
	}
	app.Result(c, resp)
}

type AckReq struct {
	Kind    ack_kind.Code `json:"kind"`
	MsgID   int64         `json:"msg_id"`
	OrderNo string        `json:"order_no"`
}

func Ack(c *gin.Context) {
	req := new(AckReq)
	if err := c.ShouldBindJSON(req); err != nil {
		app.InvalidParams(c, err.Error())
		return
	}
	if req.Kind <= 0 || (req.MsgID <= 0 && req.OrderNo == "") {
		app.InvalidParams(c, "empty parameters")
		return
	}
	var (
		ctx      = c.Request.Context()
		rds      = rdb.Client
		memberID = auth.MemberID(c)
	)
	switch req.Kind {
	case ack_kind.Overlay, ack_kind.Center:
		// 悬浮框、首页中间通知->已读
		u := unread.NewFeed(rds, memberID)
		if err := u.Ack(ctx, strconv.FormatInt(req.MsgID, 10)); err != nil {
			zapx.ErrorCtx(ctx, "ack error", zap.Error(err))
		}
		// 通知中心->已读
		u2 := unread.NewScopeSet[int64](rds, memberID, proto.ScopeCenter)
		if err := u2.Ack(ctx, req.MsgID); err != nil {
			zapx.ErrorCtx(ctx, "ack error", zap.Error(err))
		}
		syncDB(ctx, memberID, req.MsgID)
	case ack_kind.CenterAll:
		// 通知中心->全部已读
		u := unread.NewScopeSet[int64](rds, memberID, proto.ScopeCenter)
		if err := u.Del(ctx); err != nil {
			zapx.ErrorCtx(ctx, "ack error", zap.Error(err))
		}
		syncDB(ctx, memberID, 0)
	case ack_kind.OrderTab:
		// 订单中心->已读
		u := unread.NewScopeSet[string](rds, memberID, proto.ScopeTabOrder)
		if err := u.Ack(ctx, req.OrderNo); err != nil {
			zapx.ErrorCtx(ctx, "ack error", zap.Error(err))
		}
	}
	app.Success(c)
}

func syncDB(ctx context.Context, memberID int64, msgID int64) {
	_, _ = im_rpcx.AckNotify(ctx, &contracts.AckNotifyReq{
		MemberID: memberID,
		MsgID:    msgID,
	})
}
