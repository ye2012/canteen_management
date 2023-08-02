package dto

import "github.com/canteen_management/enum"

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

type WeekMenuListReq MenuListReq

type WeekMenuInfo struct {
	WeekMenuID    uint32   `json:"week_menu_id"`
	MenuType      uint32   `json:"menu_type"`
	MenuStartDate int64    `json:"menu_start_time"`
	MenuEndDate   int64    `json:"menu_end_time"`
	MenuContent   []string `json:"menu_content"`
}

type WeekMenuListRes struct {
	MenuList []*WeekMenuInfo `json:"menu_list"`
}

type WeekMenuDetailReq struct {
	WeekMenuID uint32 `json:"week_menu_id"`
}

type WeekMenuDetailInfo struct {
	WeekMenuID    uint32      `json:"week_menu_id"`
	MenuType      uint32      `json:"menu_type"`
	MenuStartDate int64       `json:"menu_start_time"`
	MenuEndDate   int64       `json:"menu_end_time"`
	MenuList      []*MenuInfo `json:"menu_list"`
}

type WeekMenuDetailRes WeekMenuDetailInfo

type ModifyWeekMenuReq struct {
	Operate  enum.OperateType   `json:"operate"`
	WeekMenu WeekMenuDetailInfo `json:"week_menu"`
}

type GenerateStaffMenuReq struct {
	MenuType uint32 `json:"menu_type"`
	MenuDate int64  `json:"menu_date"`
}

type GenerateStaffMenuRes = MenuInfo

type GenerateWeekMenuReq struct {
	MenuType  uint32 `json:"menu_type"`
	TimeStart int64  `json:"time_start"`
}

type GenerateWeekMenuRes = WeekMenuDetailInfo

type StaffMenuListHeadReq struct {
}

type StaffMenuListHeadRes []*TableColumnInfo

type StaffMenuListDataReq struct {
	TimeStart int64 `json:"time_start"`
	TimeEnd   int64 `json:"time_end"`
}

type StaffMenuListDataRes []*TableRowInfo

type StaffMenuDetailHeadReq struct {
}

type StaffMenuDetailHeadRes []*TableColumnInfo

type StaffMenuDetailDataReq struct {
	StaffMenuID uint32 `json:"staff_menu_id"`
}

type StaffMenuDetailDataRes []*TableRowInfo

type ModifyStaffMenuDetailReq struct {
	Operate    enum.OperateType `json:"operate"`
	MenuDetail []*TableRowInfo  `json:"menu_detail"`
}

type MenuTypeListHeadReq struct {
}
type MenuTypeListHeadRes []*TableColumnInfo

type MenuTypeListDataReq struct {
}

type MenuTypeListDataRes []*TableRowInfo

type MenuTypeDetailHeadReq struct {
	MenuTypeID uint32 `json:"menu_type_id"`
}

type MenuTypeDetailHeadRes []*TableColumnInfo

type MenuTypeDetailDataReq struct {
	MenuTypeID uint32 `json:"menu_type_id"`
}

type MenuTypeDetailDataRes []*TableRowInfo

type ModifyMenuTypeReq struct {
	Operate  enum.OperateType `json:"operate"`
	MenuType []*TableRowInfo  `json:"menu_type"`
}
