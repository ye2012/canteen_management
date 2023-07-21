package enum

type PurchaseOrderStatus = int8

const (
	PurchaseOrderStatusAll                     = -1
	PurchaseOrderNew       PurchaseOrderStatus = iota
	PurchaseOrderAudited
	PurchaseOrderAccept
	PurchaseOrderReceive
	PurchaseOrderFinish
)
