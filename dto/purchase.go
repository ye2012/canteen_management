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
	SupplierList []*SupplierInfo `json:"supplier_list"`
}

type ModifySupplierReq struct {
	Operate  enum.OperateType `json:"operate"`
	Supplier *SupplierInfo    `json:"supplier"`
}

type BindSupplierReq struct {
	SupplierID uint32 `json:"supplier_id"`
	OpenID     string `json:"open_id"`
}

type GoodsPurchaseInfo struct {
	GoodsId       int64   `json:"goods_id"`
	Name          string  `json:"name"`
	GoodsTypeID   uint32  `json:"goods_type_id"`
	ExpectAmount  float64 `json:"expect_amount"`
	ReceiveAmount float64 `json:"receive_amount"`
	Discount      float64 `json:"discount"`
	DealPrice     float64 `json:"deal_price"`
}

type PurchaseOrderInfo struct {
	ID          uint32               `json:"id"`
	Supplier    uint32               `json:"supplier"`
	SignPicture []string             `json:"sign_picture"`
	Status      uint8                `json:"status"`
	GoodsList   []*GoodsPurchaseInfo `json:"goods_list"`
}

type PurchaseOrderListReq struct {
}

type ApplyPurchaseReq struct {
	GoodsList []*GoodsPurchaseInfo `json:"goods_list"`
}

type ConfirmPurchaseReq struct {
	PurchaseID uint32 `json:"purchase_id"`
}

type ReceivePurchaseReq struct {
	PurchaseID uint32               `json:"purchase_id"`
	GoodsList  []*GoodsPurchaseInfo `json:"goods_list"`
}
