package enum

type OperateType = uint8

const (
	OperateTypeAdd = iota + 1
	OperateTypeModify
	OperateTypeDel
)
