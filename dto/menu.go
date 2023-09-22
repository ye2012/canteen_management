package dto

import "github.com/canteen_management/enum"

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

type WeekMenuListReq struct {
	PaginationReq
	MenuType  uint32 `json:"menu_type"`
	TimeStart int64  `json:"time_start"`
	TimeEnd   int64  `json:"time_end"`
}

type WeekMenuInfo struct {
	WeekMenuID    uint32   `json:"week_menu_id"`
	MenuType      uint32   `json:"menu_type"`
	MenuStartDate int64    `json:"menu_start_time"`
	MenuEndDate   int64    `json:"menu_end_time"`
	MenuContent   []string `json:"menu_content"`
}

type WeekMenuListRes struct {
	PaginationRes
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

type WeekMenuListHeadReq = WeekMenuListReq

type WeekMenuDetailTableReq = WeekMenuDetailReq
type WeekMenuDetailTableRes struct {
	Head []*TableColumnInfo `json:"head"`
	Data []TableRowInfo     `json:"data"`
}

type ModifyWeekMenuDetailReq struct {
	Operate      enum.OperateType `json:"operate"`
	WeekMenuID   uint32           `json:"week_menu_id"`
	MenuTypeID   uint32           `json:"menu_type_id"`
	MenuDate     int64            `json:"menu_date"`
	WeekMenuRows []*TableRowInfo  `json:"week_menu_rows"`
}

type ModifyWeekMenuReq struct {
	Operate  enum.OperateType   `json:"operate"`
	WeekMenu WeekMenuDetailInfo `json:"week_menu"`
}

type GenerateStaffMenuReq struct {
	MenuType uint32 `json:"menu_type"`
}

type GenerateStaffMenuRes struct {
	Head []*TableColumnInfo `json:"head"`
	Data []TableRowInfo     `json:"data"`
}

type GenerateWeekMenuReq struct {
	MenuType  uint32 `json:"menu_type"`
	TimeStart int64  `json:"time_start"`
}

type GenerateWeekMenuRes = GenerateStaffMenuRes

type StaffMenuListHeadReq struct {
}

type StaffMenuListHeadRes []*TableColumnInfo

type StaffMenuListDataReq struct {
	PaginationReq
	TimeStart int64 `json:"time_start"`
	TimeEnd   int64 `json:"time_end"`
}

type StaffMenuListDataRes struct {
	PaginationRes
	TableRowList []*TableRowInfo `json:"table_row_list"`
}

type StaffMenuDetailHeadReq struct {
}

type StaffMenuDetailHeadRes []*TableColumnInfo

type StaffMenuDetailDataReq struct {
	StaffMenuID uint32 `json:"staff_menu_id"`
}

type StaffMenuDetailDataRes []*TableRowInfo

type ModifyStaffMenuDetailReq struct {
	Operate       enum.OperateType `json:"operate"`
	StaffMenuID   uint32           `json:"staff_menu_id"`
	MenuTypeID    uint32           `json:"menu_type_id"`
	MenuDate      int64            `json:"menu_date"`
	StaffMenuRows []*TableRowInfo  `json:"staff_menu_rows"`
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
	Operate      enum.OperateType `json:"operate"`
	MenuTypeID   uint32           `json:"menu_type_id"`
	MenuTypeName string           `json:"menu_type_name"`
	MenuTypeRows []*TableRowInfo  `json:"menu_type_rows"`
}
