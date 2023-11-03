package dto

import (
	"fmt"

	"github.com/canteen_management/enum"
)

type SupplierListReq struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type SupplierInfo struct {
	SupplierID       uint32 `json:"supplier_id"`
	Name             string `json:"name"`
	PhoneNumber      string `json:"phone_number"`
	IDNumber         string `json:"id_number"`
	Location         string `json:"location"`
	ValidityDeadline int64  `json:"validity_deadline"`
	OpenID           string `json:"open_id"`
}

type SupplierListRes struct {
	SupplierList     []*SupplierInfo `json:"supplier_list"`
	LastValidityTime int64           `json:"last_validity_time"`
}

type ModifySupplierReq struct {
	Operate  enum.OperateType `json:"operate"`
	Supplier *SupplierInfo    `json:"supplier"`
}

type BindSupplierReq struct {
	SupplierID uint32 `json:"supplier_id"`
	OpenID     string `json:"open_id"`
}

type RenewSupplierReq struct {
	SupplierID uint32 `json:"supplier_id"`
	EndTime    int64  `json:"end_time"`
}

type PurchaseGoodsInfo struct {
	PurchaseGoodsBase
	ReceiveNumber float64 `json:"receive_number"`
	Price         float64 `json:"price"`
}

type PurchaseOrderInfo struct {
	ID            uint32               `json:"id"`
	Supplier      uint32               `json:"supplier"`
	SupplierName  string               `json:"supplier_name"`
	GoodsList     []*PurchaseGoodsInfo `json:"goods_list"`
	TotalAmount   float64              `json:"total_amount"`
	PaymentAmount float64              `json:"payment_amount"`
	Status        uint8                `json:"status"`
	Creator       string               `json:"creator"`
	CreateTime    int64                `json:"create_time"`
	ReceiveTime   int64                `json:"receive_time"`
	Receiver      string               `json:"receiver"`
	SignPicture   string               `json:"sign_picture"`
}

type PurchaseListReq struct {
	PaginationReq
	Status     int8   `json:"status"`
	Uid        uint32 `json:"uid"`
	PurchaseID uint32 `json:"purchase_id"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
}

type PurchaseListRes struct {
	PaginationRes
	PurchaseList []*PurchaseOrderInfo `json:"purchase_list"`
}

type ApplyPurchaseReq struct {
	Uid       uint32               `json:"uid"`
	GoodsList []*PurchaseGoodsInfo `json:"goods_list"`
}

type ReviewPurchaseReq struct {
	PurchaseID uint32 `json:"purchase_id"`
}

type ConfirmPurchaseReq struct {
	PurchaseID uint32 `json:"purchase_id"`
}

type ReceivePurchaseReq struct {
	PurchaseID  uint32               `json:"purchase_id"`
	Uid         uint32               `json:"uid"`
	GoodsList   []*PurchaseGoodsInfo `json:"goods_list"`
	SignPicture string               `json:"sign_picture"`
}

func (rpr *ReceivePurchaseReq) CheckParams() error {
	if rpr.SignPicture == "" {
		return fmt.Errorf("请完成签名")
	}
	return nil
}
