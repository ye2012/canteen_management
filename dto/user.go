package dto

import "fmt"

type CanteenLoginReq struct {
	Code string `json:"code"`
}

type CanteenLoginRes struct {
	Uid         uint32  `json:"uid"`
	OpenID      string  `json:"open_id"`
	PhoneNumber string  `json:"phone_number"`
	Discount    float64 `json:"discount"`
	ExtraPay    float64 `json:"extra_pay"`
	Role        uint32  `json:"role"`
}

type BindPhoneNumberReq struct {
	Uid         uint32 `json:"uid"`
	PhoneNumber string `json:"phone_number"`
}

func (bpn *BindPhoneNumberReq) CheckParams() error {
	if len(bpn.PhoneNumber) < 11 {
		return fmt.Errorf("请输入正确的手机号格式")
	}
	return nil
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
