package enum

type RoleType = uint8

const (
	RoleMin RoleType = iota
	RoleAdmin
	RoleDeliver
	RolePurchaser
	RoleSupplier
	RoleReviewer
	RoleReceiver
	RoleApplier
	RoleMax
)
