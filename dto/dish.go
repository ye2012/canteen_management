package dto

import (
	"github.com/canteen_management/enum"
)

type DishTypeListReq struct {
}

type DishTypeInfo struct {
	DishTypeID   uint32 `json:"dish_type_id"`
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
	DishID   uint32  `json:"dish_id"`
	DishName string  `json:"dish_name"`
	DishType uint32  `json:"dish_type"`
	Material string  `json:"material"`
	Price    float64 `json:"price"`
}

type DishListRes struct {
	DishList []*DishInfo `json:"dish_list"`
}

type ModifyDishReq struct {
	Operate  enum.OperateType `json:"operate"`
	DishInfo DishInfo         `json:"dish_info"`
}

type MenuListReq struct {
	MenuType  uint32 `json:"menu_type"`
	TimeStart int64  `json:"time_start"`
	TimeEnd   int64  `json:"time_end"`
}

type MealInfo struct {
	MealName string      `json:"meal_name"`
	MealType uint8       `json:"meal_type"`
	DishList []*DishInfo `json:"dish_list"`
}

type MenuInfo struct {
	MenuID   uint32      `json:"menu_id"`
	MenuType uint32      `json:"menu_type"`
	MenuDate int64       `json:"menu_time"`
	MealList []*MealInfo `json:"meal_list"`
}

type MenuListRes struct {
	MenuList []*MenuInfo `json:"menu_list"`
}

type ModifyMenuReq struct {
	Operate enum.OperateType `json:"operate"`
	Menu    MenuInfo         `json:"menu"`
}

type MenuTypeListReq struct {
}

type MenuConfigInfo struct {
	MealName      string           `json:"meal_name"`
	MealType      uint8            `json:"meal_type"`
	DishNumberMap map[uint32]int32 `json:"dish_number_map"`
}

type MenuTypeInfo struct {
	MenuTypeID uint32            `json:"menu_type_id"`
	Name       string            `json:"name"`
	MenuConfig []*MenuConfigInfo `json:"menu_config"`
}

type MenuTypeListRes struct {
	TypeList []*MenuTypeInfo `json:"type_list"`
}

type ModifyMenuTypeReq struct {
	Operate  enum.OperateType `json:"operate"`
	MenuType MenuTypeInfo     `json:"menu_type"`
}

type GenerateMenuReq struct {
	MenuType  uint32 `json:"menu_type"`
	TimeStart int64  `json:"time_start"`
	TimeEnd   int64  `json:"time_end"`
}

type GenerateMenuRes struct {
	MenuList []*MenuInfo `json:"menu_list"`
}
