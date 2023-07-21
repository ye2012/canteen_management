package enum

type RoleType = uint8

const (
	RolePurchaser = iota + 1
	RoleReviewer
	RoleSupplier
	RoleReceiver
	RoleApplier
)
