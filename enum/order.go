package enum

type PayOrderStatus = int8

const (
	PayOrderNew = iota
	PayOrderFinish
	PayOrderTimeOut
	PayOrderCancel
)

type OrderStatus = int8

const (
	OrderNew = iota
	OrderPaid
	OrderCancel
	OrderReady
	OrderFinish
)

type PurchaseStatus = int8

const (
	PurchaseNew PurchaseStatus = iota
	PurchaseReviewed
	PurchaseAccept
	PurchaseReceived
	PurchaseFinish
)

var buildingMap = map[uint32]string{
	1: "A座",
	2: "B座",
}

func GetBuildingName(buildingID uint32) string {
	name, ok := buildingMap[buildingID]
	if ok {
		return name
	}
	return ""
}

type PayMethod = uint8

const (
	PayMethodWeChat PayMethod = iota
	PayMethodCash
)

type OutboundStatus = int8

const (
	OutboundNew OutboundStatus = iota
	OutboundReviewed
)

type InventoryOrderStatus = int8

const (
	InventoryOrderNew InventoryOrderStatus = iota
	InventoryOrderFinish
	InventoryOrderConfirmed
	InventoryOrderReviewed
)

type InventoryStatus = int8

const (
	InventoryNew InventoryStatus = iota
	InventoryMatch
	InventoryNeedFix
)
