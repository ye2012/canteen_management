package model

type OutboundOrderDetail struct {
	ID         uint32  `json:"id"`
	OutboundID uint32  `json:"outbound_id"`
	GoodsID    uint32  `json:"goods_id"`
	GoodsType  uint32  `json:"goods_type"`
	OutNumber  float64 `json:"out_number"`
	Price      float64 `json:"price"`
}
