package dto

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
	Page        uint32 `json:"page"`
	PageSize    uint32 `json:"page_size"`
}

type OrderListRes = []*OrderInfo
