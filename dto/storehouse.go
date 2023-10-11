package dto

import "github.com/canteen_management/enum"

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
	GoodsList []*OutboundGoodsInfo `json:"goods_list"`
}

type OutboundListReq struct {
	PaginationReq
	Uid        uint32 `json:"uid"`
	PurchaseID uint32 `json:"outbound_id"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
}

type OutboundOrderInfo struct {
	ID          uint32               `json:"id"`
	GoodsList   []*OutboundGoodsInfo `json:"goods_list"`
	TotalAmount float64              `json:"total_amount"`
}

type OutboundListRes struct {
	PaginationRes
	OutboundList []*OutboundOrderInfo `json:"outbound_list"`
}
