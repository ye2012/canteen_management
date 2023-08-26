package dto

import "github.com/canteen_management/enum"

type OrderMenuReq struct {
}

type OrderMenuRes = []*OrderNode

type ApplyItem struct {
	DishID   uint32  `json:"dish_id"`
	DishName string  `json:"dish_name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

type ApplyOrderReq = OrderInfo

type OrderInfo struct {
	ID            string       `json:"id"`
	UnionID       string       `json:"union_id"`
	OrderID       string       `json:"order_id"`
	OrderNo       string       `json:"order_no"`
	Address       string       `json:"address"`
	PickUpMethod  uint8        `json:"pick_up_method"`
	TotalAmount   float64      `json:"total_amount"`
	PaymentAmount float64      `json:"payment_amount"`
	OrderItems    []*ApplyItem `json:"order_items"`
	CreateTime    int64        `json:"create_time"`
	OrderStatus   uint8        `json:"order_status"`
}

type ApplyOrderRes struct {
	Order     *OrderInfo `json:"order"`
	PrepareID string     `json:"prepare_id"`
}

type OrderListReq struct {
	OrderStatus int8   `json:"order_status"`
	Uid         uint32 `json:"uid"`
	OrderID     uint32 `json:"order_id"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	Page        uint32 `json:"page"`
	PageSize    uint32 `json:"page_size"`
}

type OrderListRes = []*OrderInfo

type OrderUserListReq struct {
	PhoneNumber   string `json:"phone_number"`
	DiscountLevel int32  `json:"discount_level"`
	Page          uint32 `json:"page"`
	PageSize      uint32 `json:"page_size"`
}

type OrderUserInfo struct {
	ID            uint32 `json:"id"`
	PhoneNumber   string `json:"phone_number"`
	DiscountLevel int32  `json:"discount_level"`
}

type OrderUserListRes struct {
	UserList  []*OrderUserInfo `json:"user_list"`
	TotalPage uint32           `json:"total_page"`
	Page      uint32           `json:"page"`
	PageSize  uint32           `json:"page_size"`
}

type ModifyOrderUserReq struct {
	Operate  enum.OperateType `json:"operate"`
	UserList []*OrderUserInfo `json:"user_list"`
}

type OrderDiscountListReq struct {
}

type OrderDiscountInfo struct {
	ID                uint32  `json:"id"`
	DiscountTypeName  string  `json:"discount_type_name"`
	BreakfastDiscount float64 `json:"breakfast_discount"`
	LunchDiscount     float64 `json:"lunch_discount"`
	DinnerDiscount    float64 `json:"dinner_discount"`
}

type OrderDiscountListRes struct {
	DiscountList []*OrderDiscountInfo `json:"discount_list"`
}

type ModifyOrderDiscountReq struct {
	Operate      enum.OperateType   `json:"operate"`
	DiscountInfo *OrderDiscountInfo `json:"discount_info"`
}
