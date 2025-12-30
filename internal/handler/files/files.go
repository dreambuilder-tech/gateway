package files

import (
	"gateway/internal/common/auth"
	"gateway/internal/common/file"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/upload_scene"
	"wallet/common-lib/rpcx/contracts"
	"wallet/common-lib/rpcx/member_rpcx"
	"wallet/common-lib/utils/s3x"
	"wallet/common-lib/zapx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadFileResp struct {
	FileID int64  `json:"file_id"`
	URL    string `json:"url"`
}

func Upload(c *gin.Context) {
	s := c.PostForm("scene")
	if s == "" {
		app.InvalidParams(c, "empty scene")
		return
	}
	fileUpload(c, upload_scene.Code(s))
}

func fileUpload(c *gin.Context, scene upload_scene.Code) {
	r, err := file.UpLoadImage(c)
	if err != nil {
		zapx.ErrorCtx(c.Request.Context(), "upload image error", zap.Error(err))
		return
	}
	var URL string
	switch scene {
	case upload_scene.TradeHistoryCapture:
		URL, err = s3x.UploadTradeHistoryCapture(r.FullName, r.Bytes)
	case upload_scene.ChatImage:
		URL, err = s3x.UploadChatImage(r.FullName, r.Bytes)
	case upload_scene.TradeProofCapture:
		URL, err = s3x.UploadTradeProof(r.FullName, r.Bytes)
	case upload_scene.IDImage:
		URL, err = s3x.UploadIdImage(r.FullName, r.Bytes)
	case upload_scene.PaymentQRCode:
		URL, err = s3x.UploadPaymentQRCode(r.FullName, r.Bytes)
	default:
		app.InvalidParams(c, "unsupported scene")
		return
	}
	if err != nil {
		zapx.ErrorCtx(c.Request.Context(), "upload to s3 bucket error", zap.Error(err))
		app.InternalError(c, "upload to storage failed")
		return
	}
	if URL == "" {
		app.InternalError(c, "upload s3 failed. empty URL")
		return
	}
	resp, err := member_rpcx.UploadFile(c.Request.Context(), &contracts.UploadFileReq{
		MemberID: auth.MemberID(c),
		Scene:    scene,
		S3URL:    URL,
		FileSize: len(r.Bytes),
		Mime:     r.Mime,
	})
	if err != nil {
		zapx.ErrorCtx(c.Request.Context(), "upload file error", zap.Error(err))
		app.InternalError(c, "upload file failed")
		return
	}
	app.Result(c, &UploadFileResp{
		FileID: resp.FileID,
		URL:    URL,
	})
}
