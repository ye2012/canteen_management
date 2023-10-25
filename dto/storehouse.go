package dto

import (
	"github.com/canteen_management/enum"
)

type StoreTypeListReq struct {
}

type StoreTypeInfo struct {
	StoreTypeID   uint32 `json:"store_type_id"`
	StoreTypeName string `json:"store_type_name"`
}

type StoreTypeListRes struct {
	StoreTypeList []*StoreTypeInfo `json:"store_type_list"`
}

type ModifyStoreTypeReq struct {
	Operate   enum.OperateType `json:"operate"`
	StoreType *StoreTypeInfo   `json:"store_type"`
}

type StoreListReq struct {
	StoreType uint32 `json:"store_type"`
}

type StoreGoodsInfo struct {
}

type StoreListRes struct {
	GoodsList []*GoodsInfo `json:"goods_list"`
}

type ResetStoreGoodsQuantity struct {
	Goods *GoodsTypeInfo `json:"goods"`
}

type OutboundGoodsInfo struct {
	PurchaseGoodsBase
}

type ApplyOutboundReq struct {
	Uid       uint32               `json:"uid"`
	GoodsList []*OutboundGoodsInfo `json:"goods_list"`
}

type ReviewOutboundReq struct {
	OutboundID uint32 `json:"outbound_id"`
}

type OutboundListReq struct {
	PaginationReq
	Uid        uint32 `json:"uid"`
	OutboundID uint32 `json:"outbound_id"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
}

type OutboundOrderInfo struct {
	ID               uint32               `json:"id"`
	TotalGoodsNumber int32                `json:"total_goods_number"`
	TotalGoodsType   int32                `json:"total_goods_type"`
	TotalWeight      float64              `json:"total_weight"`
	GoodsList        []*OutboundGoodsInfo `json:"goods_list"`
	TotalAmount      float64              `json:"total_amount"`
	OutboundTime     int64                `json:"outbound_time"`
	Sender           string               `json:"sender"`
}

type OutboundListRes struct {
	PaginationRes
	OutboundList []*OutboundOrderInfo `json:"outbound_list"`
}

type InventoryListReq struct {
	PaginationReq
	Uid         uint32 `json:"uid"`
	InventoryID uint32 `json:"inventory_id"`
	Status      int8   `json:"status"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
}

type InventoryGoodsNode struct {
	PurchaseGoodsBase
	BatchSize  float64               `json:"batch_size"`
	BatchUnit  string                `json:"batch_unit"`
	RealNumber float64               `json:"real_number"`
	Tag        string                `json:"tag"`
	Status     int8                  `json:"status"`
	Children   []*InventoryGoodsNode `json:"children,omitempty"`
}

type InventoryOrderInfo struct {
	ID              uint32                `json:"id"`
	TotalCount      int32                 `json:"total_count"`
	ExceptionCount  int32                 `json:"exception_count"`
	ExceptionAmount float64               `json:"exception_amount"`
	Status          int8                  `json:"status"`
	GoodsList       []*InventoryGoodsNode `json:"goods_list"`
	StartTime       int64                 `json:"start_time"`
	EndTime         int64                 `json:"end_time"`
	Creator         string                `json:"creator"`
	Partner         string                `json:"partner"`
}

type InventoryListRes struct {
	PaginationRes
	InventoryList []*InventoryOrderInfo `json:"inventory_list"`
}

type InventoryReq struct {
	InventoryID        uint32              `json:"inventory_id"`
	InventoryGoodsInfo *InventoryGoodsNode `json:"inventory_goods_info"`
}

type ApplyInventoryReq struct {
	Uid         uint32 `json:"uid"`
	InventoryID uint32 `json:"inventory_id"`
}

type ConfirmInventoryReq struct {
	Uid         uint32 `json:"uid"`
	InventoryID uint32 `json:"inventory_id"`
}

type ReviewInventoryReq struct {
	InventoryID uint32 `json:"inventory_id"`
}

type StartInventoryReq struct {
	Uid uint32 `json:"uid"`
	New bool   `json:"new"`
}

type StartInventoryRes struct {
	InventoryOrder *InventoryOrderInfo `json:"inventory_order"`
}

type GoodsHistoryReq struct {
	PaginationReq
	GoodsID    uint32 `json:"goods_id"`
	ChangeType uint32 `json:"change_type"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
}

type GoodsHistoryInfo struct {
	ID             uint32  `json:"id"`
	GoodsID        uint32  `json:"goods_id"`
	ChangeQuantity float64 `json:"change_quantity"`
	BeforeQuantity float64 `json:"before_quantity"`
	AfterQuantity  float64 `json:"after_quantity"`
	ChangeType     uint32  `json:"change_type"`
	RefID          uint32  `json:"ref_id"`
	CreateAt       int64   `json:"created_at"`
}

type GoodsHistoryRes struct {
	PaginationRes
	History []*GoodsHistoryInfo `json:"history"`
}
