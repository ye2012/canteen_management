package dto

import (
	"fmt"

	"github.com/canteen_management/enum"
)

type AdminLoginReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (alr *AdminLoginReq) CheckParams() error {
	if len(alr.Password) > 20 {
		return fmt.Errorf("密码长度不合法")
	}
	return nil
}

type RouterNode struct {
	Name     string        `json:"name"`
	Path     string        `json:"path,omitempty"`
	Children []*RouterNode `json:"children,omitempty"`
}

type AdminLoginRes struct {
	Router []*RouterNode `json:"router"`
	Token  string        `json:"token"`  //登录成功返回token
	Expire int64         `json:"expire"` //token过期时间
}

type CanteenLoginReq struct {
	Code string `json:"code"`
}

type CanteenLoginRes struct {
	Uid          uint32   `json:"uid"`
	OpenID       string   `json:"open_id"`
	PhoneNumber  string   `json:"phone_number"`
	Discount     float64  `json:"discount"`
	DiscountLeft float64  `json:"discount_left"`
	ExtraPay     float64  `json:"extra_pay"`
	RoleList     []uint32 `json:"role_list"`
	Token        string   `json:"token"`  //登录成功返回token
	Expire       int64    `json:"expire"` //token过期时间
}

type CanteenUserCenterReq struct {
	Uid uint32 `json:"uid"`
}

type CanteenUserCenterRes struct {
	PhoneNumber  string   `json:"phone_number"`
	Discount     float64  `json:"discount"`
	DiscountLeft float64  `json:"discount_left"`
	ExtraPay     float64  `json:"extra_pay"`
	RoleList     []uint32 `json:"role_list"`
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
	Code string `json:"code"`
}

type KitchenLoginRes struct {
	Uid         uint32   `json:"uid"`
	OpenID      string   `json:"open_id"`
	PhoneNumber string   `json:"phone_number"`
	RoleList    []uint32 `json:"role_list"`
	Token       string   `json:"token"`  //登录成功返回token
	Expire      int64    `json:"expire"` //token过期时间
}

type KitchenUserCenterReq struct {
	Uid uint32 `json:"uid"`
}

type KitchenUserCenterRes struct {
	PhoneNumber string   `json:"phone_number"`
	RoleList    []uint32 `json:"role_list"`
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

type RouterTypeListReq struct {
}

type RouterTypeInfo struct {
	RouterTypeID   uint32 `json:"router_type_id"`
	RouterTypeName string `json:"router_type_name"`
	SortID         uint32 `json:"sort_id"`
}

type RouterTypeListRes struct {
	RouterTypeList []*RouterTypeInfo `json:"router_type_list"`
}

type ModifyRouterTypeReq struct {
	Operate    enum.OperateType `json:"operate"`
	RouterType *RouterTypeInfo  `json:"router_type"`
}

type RouterListReq struct {
	RouterType uint32 `json:"router_type"`
}

type RouterInfo struct {
	RouterID     uint32   `json:"router_id"`
	RouterType   uint32   `json:"router_type"`
	RouterName   string   `json:"router_name"`
	RouterPath   string   `json:"router_path"`
	RouterSortID uint32   `json:"router_sort_id"`
	RoleList     []uint32 `json:"role_list"`
}

func (ri *RouterInfo) CheckParams() error {
	if ri.RouterPath == "" {
		return fmt.Errorf("路由Path必须设置")
	}
	if ri.RouterName == "" {
		return fmt.Errorf("路由名字必须设置")
	}
	if ri.RouterID == 0 {
		return fmt.Errorf("路由类型必须设置")
	}
	return nil
}

type RouterListRes struct {
	RouterList []*RouterInfo `json:"router_list"`
}

type ModifyRouterReq struct {
	Operate enum.OperateType `json:"operate"`
	Router  *RouterInfo      `json:"router"`
}
