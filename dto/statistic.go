package dto

type DashboardRes struct {
	TotalOrderInfo
	DayOrderInfo
}

type DayOrderInfo struct {
	DayPayAmount      float64 `json:"day_pay_amount"`
	DayOrderCount     int32   `json:"day_order_count"`
	DaySuccessCount   int32   `json:"day_success_count"`
	DayBreakfastCount int32   `json:"day_breakfast_count"`
	DayLunchCount     int32   `json:"day_lunch_count"`
	DayDinnerCount    int32   `json:"day_dinner_count"`
}

type TotalOrderInfo struct {
	TotalPayAmount    float64 `json:"total_pay_amount"`
	TotalOrderCount   int32   `json:"total_order_count"`
	TotalSuccessCount int32   `json:"total_success_count"`
}
