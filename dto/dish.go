package dto

import (
	"github.com/canteen_management/enum"
)

type DishTypeListReq struct {
	MasterTypeID     uint32 `json:"master_type_id"`
	IncludeMaserType bool   `json:"include_maser_type"`
}

type DishTypeInfo struct {
	DishTypeID   uint32 `json:"dish_type_id"`
	MasterTypeID uint32 `json:"master_type_id"`
	DishTypeName string `json:"dish_type_name"`
}

type DishTypeListRes struct {
	List []*DishTypeInfo `json:"list"`
}

type ModifyDishTypeReq struct {
	Operate  enum.OperateType `json:"operate"`
	TypeInfo DishTypeInfo     `json:"dish_type_info"`
}

type DishListReq struct {
	DishType uint32 `json:"dish_type"`
}

type DishInfo struct {
	DishID         uint32  `json:"dish_id"`
	DishName       string  `json:"dish_name"`
	DishTypeID     uint32  `json:"dish_type_id"`
	DishTypeName   string  `json:"dish_type_name"`
	MasterTypeName string  `json:"master_type_name"`
	Material       string  `json:"material"`
	Price          float64 `json:"price"`
}

type DishListRes struct {
	DishList []*DishInfo `json:"dish_list"`
}

type ModifyDishReq struct {
	Operate  enum.OperateType `json:"operate"`
	DishInfo DishInfo         `json:"dish_info"`
}
