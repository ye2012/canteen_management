package enum

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
