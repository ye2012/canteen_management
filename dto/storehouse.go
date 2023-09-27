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

type GoodsOutboundInfo struct {
	GoodsId       int64   `json:"goods_id"`
	Name          string  `json:"name"`
	GoodsTypeID   uint32  `json:"goods_type_id"`
	ExpectAmount  float64 `json:"expect_amount"`
	ReceiveAmount float64 `json:"receive_amount"`
	Discount      float64 `json:"discount"`
	DealPrice     float64 `json:"deal_price"`
}

type ApplyOutboundReq struct {
	GoodsList []*GoodsOutboundInfo `json:"goods_list"`
}
