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

type ApplyPayOrderReq = PayOrderInfo

type PayOrderInfo struct {
	ID            uint32       `json:"id"`
	OrderList     []*OrderInfo `json:"order_list"`
	Address       string       `json:"address"`
	TotalAmount   float64      `json:"total_amount"`
	PaymentAmount float64      `json:"payment_amount"`
}

type OrderInfo struct {
	ID            string       `json:"id"`
	PayOrderID    uint32       `json:"pay_order_id"`
	UserPhone     string       `json:"user_phone"`
	OrderID       string       `json:"order_id"`
	OrderNo       string       `json:"order_no"`
	Address       string       `json:"address"`
	TotalAmount   float64      `json:"total_amount"`
	PaymentAmount float64      `json:"payment_amount"`
	OrderItems    []*ApplyItem `json:"order_items"`
	CreateTime    int64        `json:"create_time"`
	OrderStatus   uint8        `json:"order_status"`
}

type ApplyOrderRes struct {
	PayOrderInfo *PayOrderInfo `json:"pay_order_info"`
	PrepareID    string        `json:"prepare_id"`
}

type PaySuccessReq struct {
	PayOrderID uint32 `json:"pay_order_id"`
}

type PayOrderListReq struct {
	PaginationReq
	OrderStatus int8   `json:"order_status"`
	Uid         uint32 `json:"uid"`
}

type PayOrderListRes struct {
	PaginationRes
	OrderList []*PayOrderInfo `json:"order_list"`
}

type OrderListReq struct {
	PaginationReq
	OrderStatus int8   `json:"order_status"`
	Uid         uint32 `json:"uid"`
	OrderID     uint32 `json:"order_id"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
}

type OrderListRes struct {
	PaginationRes
	OrderList []*OrderInfo `json:"order_list"`
}

type OrderUserListReq struct {
	PaginationReq
	PhoneNumber   string `json:"phone_number"`
	DiscountLevel int32  `json:"discount_level"`
}

type OrderUserInfo struct {
	ID            uint32 `json:"id"`
	UnionID       string `json:"union_id"`
	PhoneNumber   string `json:"phone_number"`
	DiscountLevel uint8  `json:"discount_level"`
}

type OrderUserListRes struct {
	PaginationRes
	UserList []*OrderUserInfo `json:"user_list"`
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
