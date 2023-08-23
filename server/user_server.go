package server

import (
	"github.com/canteen_management/dto"
	"github.com/gin-gonic/gin"
)

type UserServer struct {
}

func NewUserServer() *UserServer {
	return &UserServer{}
}

func (us *UserServer) RequestCanteenLogin(ctx *gin.Context, rawReq interface{}, res *dto.Response) {

}

func (us *UserServer) RequestKitchenLogin(ctx *gin.Context, rawReq interface{}, res *dto.Response) {

}
