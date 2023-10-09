package dto

import (
	"fmt"

	"github.com/canteen_management/enum"
)

type GoodsTypeListReq struct {
	PaginationReq
}

type GoodsTypeInfo struct {
	GoodsTypeID   uint32  `json:"goods_type_id"`
	GoodsTypeName string  `json:"goods_type_name"`
	Discount      float64 `json:"discount"`
}

type GoodsTypeListRes struct {
	PaginationRes
	GoodsTypeList []*GoodsTypeInfo `json:"goods_type_list"`
}

type ModifyGoodsTypeReq struct {
	Operate   enum.OperateType `json:"operate"`
	GoodsType *GoodsTypeInfo   `json:"goods_type"`
}

type GoodsListReq struct {
	PaginationReq
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
	PaginationRes
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

type GoodsPriceListReq struct {
	PaginationReq
	GoodsTypeID uint32 `json:"goods_type_id"`
	StoreTypeID uint32 `json:"store_type_id"`
}

type GoodsPriceInfo struct {
	GoodsID      uint32            `json:"goods_id"`
	GoodsName    string            `json:"goods_name"`
	PriceList    map[uint8]float64 `json:"price_list"`
	AveragePrice float64           `json:"average_price"`
}

type GoodsPriceListRes struct {
	PaginationRes
	GoodsPriceList []*GoodsPriceInfo `json:"goods_price_list"`
}

type ModifyGoodsPriceReq struct {
	GoodsID  uint32            `json:"goods_id"`
	PriceMap map[uint8]float64 `json:"price_map"`
}

type GoodsNodeListReq struct {
	CartType uint8  `json:"cart_type"`
	Uid      uint32 `json:"uid"`
}

func (gnl *GoodsNodeListReq) CheckParams() error {
	if gnl.CartType == 0 || gnl.CartType >= enum.CartTypeMax {
		return fmt.Errorf("CartType错误")
	}
	return nil
}

type GoodsNodeListRes struct {
	GoodsList  []*GoodsNode       `json:"goods_list"`
	GoodsMap   map[string]float64 `json:"goods_map"`
	TotalCost  float64            `json:"total_cost"`
	TotalGoods float64            `json:"total_goods"`
}
