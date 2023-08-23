package dto

type CanteenLoginReq struct {
	Code string `json:"code"`
}

type CanteenLoginRes struct {
}

type KitchenLoginReq struct {
	Code     string `json:"code"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type KitchenLoginRes struct {
	Code     string `json:"code"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
