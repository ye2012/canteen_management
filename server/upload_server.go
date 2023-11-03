package server

import (
	"encoding/base64"
	"github.com/canteen_management/logger"
	"io/ioutil"
	"path/filepath"

	"github.com/canteen_management/config"
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"

	"github.com/gin-gonic/gin"
)

const (
	uploadServerLogTag = "UploadServer"
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
	if err = ctx.SaveUploadedFile(file, filename); err != nil {
		logger.Warn(uploadServerLogTag, "SaveUploadedFile Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	res.Data = config.Config.FileBaseUrl + file.Filename
}

func (us UploadServer) RequestUploadBase64(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.UploadBase64Req)
	imgData, err := base64.StdEncoding.DecodeString(req.ImgData)
	if err != nil {
		logger.Warn(uploadServerLogTag, "DecodeImgString Failed|Err:%v", err)
		res.Code = enum.ParamsError
		return
	}

	if len(imgData) > 100000 {
		res.Code = enum.ParamsError
		res.Msg = "图片大小超过限制"
		return
	}

	basePath := config.Config.FileUploadPath
	fullPath := basePath + filepath.Base(req.FileName)
	err = ioutil.WriteFile(fullPath, imgData, 0666)
	if err != nil {
		logger.Warn(uploadServerLogTag, "WriteFile Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	res.Data = config.Config.FileBaseUrl + req.FileName
}
