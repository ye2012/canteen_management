package dto

import (
	"fmt"

	"github.com/canteen_management/enum"
)

type CanteenLoginReq struct {
	Code string `json:"code"`
}

type CanteenLoginRes struct {
	Uid          uint32  `json:"uid"`
	OpenID       string  `json:"open_id"`
	PhoneNumber  string  `json:"phone_number"`
	Discount     float64 `json:"discount"`
	DiscountLeft float64 `json:"discount_left"`
	ExtraPay     float64 `json:"extra_pay"`
	Role         uint32  `json:"role"`
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
	Uid         uint32 `json:"uid"`
	OpenID      string `json:"open_id"`
	PhoneNumber string `json:"phone_number"`
	Role        uint32 `json:"role"`
}

type AdminUserListReq struct {
	PaginationReq
	RoleType uint8 `json:"role_type"`
}

type UserInfo struct {
	ID          uint32   `json:"id"`
	NickName    string   `json:"nick_name"`
	UserName    string   `json:"user_name"`
	PhoneNumber string   `json:"phone_number"`
	Password    string   `json:"password"`
	RoleList    []uint32 `json:"role_list"`
	OpenID      string   `json:"open_id"`
}

type UserListRes struct {
	PaginationRes
	UserList []*UserInfo `json:"user_list"`
}

type ModifyAdminUserReq struct {
	Operate enum.OperateType `json:"operate"`
	User    *UserInfo        `json:"user"`
}

type BindAdminReq struct {
	User *UserInfo `json:"user"`
}
