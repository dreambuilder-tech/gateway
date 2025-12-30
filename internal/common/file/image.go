package file

import (
	"errors"
	"fmt"
	"wallet/common-lib/app"
	"wallet/common-lib/consts/codex"
	"wallet/common-lib/utils/imagex"

	"github.com/gin-gonic/gin"
)

func UpLoadImage(c *gin.Context) (*imagex.Result, error) {
	fh, err := c.FormFile("file")
	if err != nil {
		app.InvalidParams(c, "missing file")
		return nil, fmt.Errorf("missing file: %v", err)
	}
	f, err := fh.Open()
	if err != nil {
		app.InvalidParams(c, "open file error")
		return nil, fmt.Errorf("open file error: %v", err)
	}
	defer f.Close()

	opt := &imagex.Options{
		MaxFileBytes:     10 << 20, // 10M
		ResizeAboveBytes: 2 << 20,  // 2M
		AllowedMimes:     []string{"image/jpeg", "image/png", "image/webp"},
	}
	res, err := imagex.Process(f, opt)
	if err != nil {
		switch {
		case errors.Is(err, imagex.ErrTooLarge) || errors.Is(err, imagex.ErrTooManyPixels):
			app.Fail(c, codex.FileTooLarge, "image too large")
			return nil, err
		case errors.Is(err, imagex.ErrBadMime), errors.Is(err, imagex.ErrNotImage), errors.Is(err, imagex.ErrBadImage):
			app.Fail(c, codex.FileWrongExtension, "unsupported image extension")
			return nil, err
		default:
			app.InternalError(c, "upload image error")
			return nil, fmt.Errorf("upload image error: %v", err)
		}
	}
	return res, nil
}
