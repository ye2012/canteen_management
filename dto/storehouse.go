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
