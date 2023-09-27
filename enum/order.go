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
