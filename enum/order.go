package enum

type PayOrderStatus = int8

const (
	PayOrderNew = iota
	PayOrderFinish
	PayOrderTimeOut
)

type OrderStatus = int8

const (
	OrderNew = iota
	OrderPaid
	OrderAccept
	OrderReady
	OrderFinish
)

type PurchaseStatus = int8

const (
	PurchaseNew = iota
	PurchaseReviewed
	PurchaseAccept
	PurchaseReach
	PurchaseFinish
)
