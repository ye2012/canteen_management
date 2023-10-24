package enum

type GoodsHistoryType = uint32

const (
	GoodsInit GoodsHistoryType = iota + 1
	GoodsPurchase
	GoodsOutbound
	GoodsInventory
)
