package dto

import "github.com/canteen_management/enum"

type GoodsTypeListReq struct {
}

type GoodsTypeInfo struct {
	GoodsTypeID   uint32  `json:"goods_type_id"`
	GoodsTypeName string  `json:"goods_type_name"`
	Discount      float64 `json:"discount"`
}

type GoodsTypeListRes struct {
	GoodsTypeList []*GoodsTypeInfo `json:"goods_type_list"`
}

type ModifyGoodsTypeReq struct {
	Operate   enum.OperateType `json:"operate"`
	GoodsType *GoodsTypeInfo   `json:"goods_type"`
}

type GoodsListReq struct {
	GoodsTypeID uint32 `json:"goods_type_id"`
	StoreTypeID uint32 `json:"store_type_id"`
}

type GoodsInfo struct {
	GoodsID   uint32  `json:"goods_id"`
	GoodsName string  `json:"goods_name"`
	GoodsType uint32  `json:"goods_type"`
	StoreType uint32  `json:"store_type"`
	Picture   string  `json:"picture"`
	BatchSize float64 `json:"batch_size"`
	BatchUnit string  `json:"batch_unit"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
}

type GoodsListRes struct {
	GoodsList []*GoodsInfo `json:"goods_list"`
}

type ModifyGoodsInfoReq struct {
	Operate enum.OperateType `json:"operate"`
	Goods   *GoodsInfo       `json:"goods"`
}

type ModifyGoodsQuantityReq struct {
	Operate enum.OperateType `json:"operate"`
	Goods   *GoodsInfo       `json:"goods"`
}
