package app

import (
	"fmt"
	"net/http"
	"wallet/common-lib/consts/codex"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    codex.Code `json:"code"`
	Data    any        `json:"data"`
	Message string     `json:"message,omitempty"`
	Extend  any        `json:"extend,omitempty"`
}

type APIPageResponse struct {
	APIResponse
	Total int64 `json:"total"`
}

func Unauthorized(c *gin.Context, err string) {
	c.AbortWithStatusJSON(http.StatusOK, &APIResponse{
		Code:    codex.Unauthorized,
		Message: err,
	})
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, &APIResponse{
		Code: codex.Success,
	})
}

func SuccessData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &APIResponse{
		Code: codex.Success,
		Data: data,
	})
}

func SuccessPage(c *gin.Context, data any, total int64) {
	c.JSON(http.StatusOK, &APIPageResponse{
		APIResponse: APIResponse{
			Code: codex.Success,
			Data: data,
		},
		Total: total,
	})
}

func InvalidParams(c *gin.Context, format string, a ...any) {
	c.JSON(http.StatusOK, &APIResponse{
		Code:    codex.InvalidParams,
		Message: fmt.Sprintf(format, a...),
	})
}

func InternalError(c *gin.Context, format string, a ...any) {
	c.JSON(http.StatusOK, &APIResponse{
		Code:    codex.InternalError,
		Message: fmt.Sprintf(format, a...),
	})
}

func Failed(c *gin.Context, failCode codex.Code, format string, a ...any) {
	c.JSON(http.StatusOK, &APIResponse{
		Code:    failCode,
		Message: fmt.Sprintf(format, a...),
	})
}

func TooManyRequest(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, &APIResponse{
		Code: codex.TooManyRequest,
	})
}

func RequestExpired(c *gin.Context, err string) {
	c.AbortWithStatusJSON(http.StatusOK, &APIResponse{
		Code:    codex.RequestExpired,
		Message: err,
	})
}
