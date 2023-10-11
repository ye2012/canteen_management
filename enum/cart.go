package enum

type CartType = uint8

const (
	CartTypeOrder = iota + 1
	CartTypePurchase
	CartTypeOutbound
	CartTypeMax
)
