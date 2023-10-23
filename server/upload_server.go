package server

import (
	"path/filepath"

	"github.com/canteen_management/config"
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"

	"github.com/gin-gonic/gin"
)

type UploadServer struct {
}

func NewUploadServer() *UploadServer {
	return &UploadServer{}
}

func (us UploadServer) RequestUpload(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	file, err := ctx.FormFile("file")
	if err != nil {
		res.Code = enum.ParamsError
		return
	}

	if file.Size > 600000 {
		res.Code = enum.ParamsError
		res.Msg = "图片大小超过限制"
		return
	}

	basePath := config.Config.FileUploadPath
	filename := basePath + filepath.Base(file.Filename)
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		res.Code = enum.ParamsError
		return
	}

	res.Data = config.Config.FileBaseUrl + file.Filename
}
