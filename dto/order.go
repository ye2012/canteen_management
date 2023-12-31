package dto

import (
	"fmt"

	"github.com/canteen_management/enum"
)

type OrderMenuReq struct {
	Uid uint32 `json:"uid"`
}

type OrderMenuRes struct {
	Menu       []*OrderNode       `json:"menu"`
	GoodsMap   map[string]float64 `json:"goods_map"`
	TotalCost  float64            `json:"total_cost"`
	TotalGoods float64            `json:"total_goods"`
	CartID     uint32             `json:"cart_id"`
}

type ApplyItem struct {
	ID       string  `json:"id"`
	DishID   uint32  `json:"dish_id"`
	DishName string  `json:"dish_name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
	Picture  string  `json:"picture"`
}

type ApplyPayOrderReq PayOrderInfo

func (apo *ApplyPayOrderReq) CheckParams() error {
	if len(apo.OrderList) == 0 {
		return fmt.Errorf("请输入订单信息")
	}
	name := enum.GetBuildingName(apo.BuildingID)
	if name == "" {
		return fmt.Errorf("请输入所在楼号信息")
	}
	if apo.CartID == 0 {
		return fmt.Errorf("请添加购物车")
	}
	return nil
}

type ApplyCashOrderReq struct {
	*PayOrderInfo
	Uid uint32 `json:"uid"`
}

type ApplyStaffOrderReq struct {
	*PayOrderInfo
	Uid uint32 `json:"uid"`
}

type CancelPayOrderReq struct {
	OrderID uint32 `json:"order_id"`
}

type FinishPayOrderReq struct {
	OrderID uint32 `json:"order_id"`
}

type DeliverOrderReq struct {
	OrderID uint32 `json:"order_id"`
}

type PayOrderInfo struct {
	ID             uint32       `json:"id"`
	Uid            uint32       `json:"uid"`
	CartID         uint32       `json:"cart_id"`
	OrderList      []*OrderInfo `json:"order_list"`
	BuildingID     uint32       `json:"building_id"`
	Floor          uint32       `json:"floor"`
	Room           string       `json:"room"`
	TotalAmount    float64      `json:"total_amount"`
	PayMethod      uint8        `json:"pay_method"`
	PaymentAmount  float64      `json:"payment_amount"`
	DiscountAmount float64      `json:"discount_amount"`
	CreateTime     int64        `json:"create_time"`
	Status         uint8        `json:"status"`
}

type OrderInfo struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	PayOrderID    uint32       `json:"pay_order_id"`
	MealType      uint8        `json:"meal_type"`
	UserPhone     string       `json:"user_phone"`
	OrderID       string       `json:"order_id"`
	OrderNo       string       `json:"order_no"`
	BuildingID    uint32       `json:"building_id"`
	Floor         uint32       `json:"floor"`
	Room          string       `json:"room"`
	TotalAmount   float64      `json:"total_amount"`
	PayMethod     uint8        `json:"pay_method"`
	PaymentAmount float64      `json:"payment_amount"`
	OrderItems    []*ApplyItem `json:"order_items"`
	CreateTime    int64        `json:"create_time"`
	DeliverTime   int64        `json:"deliver_time"`
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

type FloorFilterReq struct {
	OrderDate   int64  `json:"order_date"`
	MealType    uint8  `json:"meal_type"`
	BuildingID  uint32 `json:"building_id"`
	OrderStatus int8   `json:"order_status"`
}

type FloorFilterRes struct {
	Floors []int32 `json:"floors"`
}

type OrderDishAnalysisReq struct {
	OrderDate int64  `json:"order_date"`
	MealType  uint8  `json:"meal_type"`
	DishType  uint32 `json:"dish_type"`
	DishID    uint32 `json:"dish_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type OrderDishSummaryInfo struct {
	DishID      uint32 `json:"dish_id"`
	DishName    string `json:"dish_name"`
	DishType    uint32 `json:"dish_type"`
	Quantity    int32  `json:"quantity"`
	OrderNumber int32  `json:"order_number"`
}

type OrderDishAnalysisRes struct {
	Summary []*OrderDishSummaryInfo `json:"summary"`
}

type OrderListReq struct {
	PaginationReq
	OrderStatus int8   `json:"order_status"`
	PayMethod   int8   `json:"pay_method"`
	Uid         uint32 `json:"uid"`
	OrderID     uint32 `json:"order_id"`
	BuildingID  uint32 `json:"building_id"`
	Floor       uint32 `json:"floor"`
	Room        string `json:"room"`
	MealType    uint8  `json:"meal_type"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
}

type OrderListRes struct {
	PaginationRes
	OrderList []*OrderInfo `json:"order_list"`
}

type OrderUserListReq struct {
	PaginationReq
	OpenID        string `json:"open_id"`
	PhoneNumber   string `json:"phone_number"`
	DiscountLevel int32  `json:"discount_level"`
}

type OrderUserInfo struct {
	ID            uint32 `json:"id"`
	OpenID        string `json:"open_id"`
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

type BindOrderUserReq struct {
	ID     uint32 `json:"id"`
	OpenID string `json:"open_id"`
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

type ModifyCartReq struct {
	Uid      uint32  `json:"uid"`
	CartType uint8   `json:"cart_type"`
	ItemID   string  `json:"item_id"`
	Quantity float64 `json:"quantity"`
}

type ModifyCartRes struct {
	GoodsMap   map[string]float64 `json:"goods_map"`
	TotalCost  float64            `json:"total_cost"`
	TotalGoods float64            `json:"total_goods"`
	CartID     uint32             `json:"cart_id"`
}

type GetOrderCardReq struct {
}

type OrderCardDish struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	DishID   uint32  `json:"dish_id"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type OrderCartMeal struct {
	ID    string           `json:"id"`
	Name  string           `json:"name"`
	Child []*OrderCardDish `json:"children"`
}

type GetOrderCartRes struct {
	MealList []*OrderCartMeal `json:"meal_list"`
}
