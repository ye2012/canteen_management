package server

import (
	"github.com/canteen_management/dto"
	"github.com/gin-gonic/gin"
)

type PurchaseServer struct {
}

func NewPurchaseServer() *PurchaseServer {
	return &PurchaseServer{}
}

func (ps *PurchaseServer) RequestSupplierList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {

}

func (ps *PurchaseServer) RequestModifySupplier(ctx *gin.Context, rawReq interface{}, res *dto.Response) {

}
