package dto

import "github.com/canteen_management/enum"

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
	ID            uint32  `json:"id"`
	GoodsID       uint32  `json:"goods_id"`
	Name          string  `json:"name"`
	GoodsTypeID   uint32  `json:"goods_type_id"`
	ExpectNumber  float64 `json:"expect_number"`
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
	GoodsList []*PurchaseGoodsInfo `json:"goods_list"`
}

type ReviewPurchaseReq struct {
	PurchaseID uint32 `json:"purchase_id"`
}

type ConfirmPurchaseReq struct {
	PurchaseID uint32 `json:"purchase_id"`
}

type ReceivePurchaseReq struct {
	PurchaseID uint32               `json:"purchase_id"`
	GoodsList  []*PurchaseGoodsInfo `json:"goods_list"`
}
