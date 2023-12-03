package service

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/utils"
)

const (
	userServiceLogTag = "UserService"
)

type UserService struct {
	sqlCli            *sql.DB
	adminUserModel    *model.AdminUserModel
	wxUserModel       *model.WxUserModel
	orderUserModel    *model.OrderUserModel
	supplierModel     *model.SupplierModel
	routerTypeModel   *model.RouterTypeModel
	routerDetailModel *model.RouterDetailModel
	tokenModel        *model.TokenModel
}

func NewUserService(sqlCli *sql.DB) *UserService {
	wxUserModel := model.NewWxUserModelWithDB(sqlCli)
	orderUserModel := model.NewOrderUserModel(sqlCli)
	adminUserModel := model.NewAdminUserModelWithDB(sqlCli)
	supplierModel := model.NewSupplierModelWithDB(sqlCli)
	routerTypeModel := model.NewRouterTypeModel(sqlCli)
	routerDetailModel := model.NewRouterDetailModel(sqlCli)
	tokenModel := model.NewTokenModelWithDB(sqlCli)
	return &UserService{
		sqlCli:            sqlCli,
		wxUserModel:       wxUserModel,
		orderUserModel:    orderUserModel,
		adminUserModel:    adminUserModel,
		supplierModel:     supplierModel,
		routerTypeModel:   routerTypeModel,
		routerDetailModel: routerDetailModel,
		tokenModel:        tokenModel,
	}
}

func (us *UserService) GetWxUserRole(openID string) uint32 {
	adminList, err := us.adminUserModel.GetAdminUserByCondition(" WHERE `open_id`=? ", openID)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetAdminUserByCondition Failed|OpenID:%v|Err:%v", openID, err)
		return 0
	}
	role := uint32(0)
	if len(adminList) != 0 {
		role |= adminList[0].Role
	}

	suppliers, err := us.supplierModel.GetSupplier(0, "", "", openID)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetSupplier Failed|OpenID:%v|Err:%v", openID, err)
		return role
	}
	if len(suppliers) != 0 {
		role |= uint32(1 << enum.RoleSupplier)
	}
	return role
}

func (us *UserService) GetWxUser(userID uint32) (*model.WxUser, error) {
	wxUser, err := us.wxUserModel.GetWxUserByCondition(" WHERE `id`=? ", userID)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetWxUser Failed|UserID:%v|Err:%v", userID, err)
		return nil, err
	}

	if len(wxUser) > 0 {
		return wxUser[0], nil
	}
	return nil, nil
}

func (us *UserService) GetWxUserDiscount(openID string) uint8 {
	orderUser, err := us.orderUserModel.GetOrderUserByCondition(" WHERE `open_id`=? ", openID)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetOrderUserByCondition Failed|OpenID:%v|Err:%v", openID, err)
		return 0
	}
	if len(orderUser) > 0 {
		return orderUser[0].DiscountLevel
	}
	return 0
}

func (us *UserService) WxUserLogin(openID string, role uint32) (*model.WxUser, *model.TokenDAO, error) {
	user, err := us.wxUserModel.GetWxUserByOpenID(openID)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetWxUserByOpenID Failed|OpenID:%v|Err:%v", openID, err)
		return nil, nil, err
	}

	if user != nil {
		token := us.tokenModel.LoginSuccess(user.ID, 0, role)
		return user, token, nil
	}

	wxUser := &model.WxUser{OpenID: openID}
	err = us.wxUserModel.Insert(wxUser)
	if err != nil {
		logger.Warn(userServiceLogTag, "Insert WxUser Failed|OpenID:%v|Err:%v", openID, err)
		return nil, nil, err
	}
	token := us.tokenModel.LoginSuccess(user.ID, 0, role)
	return wxUser, token, nil
}

func (us *UserService) AdminLogin(userName, password string) (*model.AdminUser, *model.TokenDAO, error) {
	condition := " WHERE `user_name` = ? "
	users, err := us.adminUserModel.GetAdminUserByCondition(condition, userName)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetAdminUserByCondition Failed|UserName:%v|Err:%v", userName, err)
		return nil, nil, err
	}

	if len(users) == 0 {
		return nil, nil, fmt.Errorf("用户不存在")
	}
	user := users[0]
	passOk := utils.ComparePass(user.Password, password)
	if passOk == false {
		return nil, nil, fmt.Errorf("密码错误")
	}

	token := us.tokenModel.LoginSuccess(0, user.ID, user.Role)
	return user, token, nil
}

func (us *UserService) CheckPhoneNumber(phoneNumber string, isExist bool) (*model.WxUser, error) {
	condition := " WHERE `phone_number` = ?  "
	users, err := us.wxUserModel.GetWxUserByCondition(condition, phoneNumber)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetWxUserByPhoneNumber Failed|phone:%v|Err:%v", phoneNumber, err)
		return nil, err
	}
	if isExist {
		if users == nil || len(users) == 0 {
			return nil, fmt.Errorf("手机号未注册，请确认手机号")
		}
		return users[0], nil
	}
	if users != nil && len(users) > 0 {
		return users[0], fmt.Errorf("手机号已绑定，请更换手机号重试")
	}
	return nil, nil
}

func (us *UserService) bindSupplier(tx *sql.Tx, openID, phoneNumber string, uid uint32) error {
	supplier, err := us.supplierModel.GetSupplier(0, "", phoneNumber, "")
	if err != nil {
		logger.Warn(userServiceLogTag, "bindSupplier GetUser Failed|phone:%v|Err:%v", phoneNumber, err)
		return err
	}
	if len(supplier) > 0 && supplier[0].OpenID == "" {
		supplier[0].OpenID = openID
		supplier[0].Uid = uid
		err = us.supplierModel.UpdateOpenIDWithTx(tx, supplier[0].ID, openID, uid)
		if err != nil {
			logger.Warn(userServiceLogTag, "bindSupplier Update Failed|SupplierUid:%v|Phone:%v|Err:%v",
				supplier[0].ID, phoneNumber, err)
			return err
		}
	}
	return nil
}

func (us *UserService) bindOrderUser(tx *sql.Tx, openID, phoneNumber string, uid uint32) error {
	orderUser, err := us.orderUserModel.GetOrderUser("", phoneNumber, 0, 1, 10)
	if err != nil {
		logger.Warn(userServiceLogTag, "bindOrderUser GetUser Failed|phone:%v|Err:%v", phoneNumber, err)
		return err
	}

	if len(orderUser) > 0 && orderUser[0].OpenID == "" {
		orderUser[0].OpenID = openID
		orderUser[0].Uid = uid
		err = us.orderUserModel.UpdateOrderUserWithTx(tx, orderUser[0], "id", "open_id", "uid")
		if err != nil {
			logger.Warn(userServiceLogTag, "bindOrderUser Update Failed|Uid:%v|OrderUser:%v|Err:%v",
				uid, orderUser[0].ID, err)
			return err
		}
	}
	return nil
}

func (us *UserService) bindAdminUser(tx *sql.Tx, openID, phoneNumber string) error {
	adminUsers, err := us.adminUserModel.GetAdminUserByCondition(" WHERE `phone_number`=? ", phoneNumber)
	if err != nil {
		logger.Warn(userServiceLogTag, "bindAdminUser GetUser Failed|phone:%v|Err:%v", phoneNumber, err)
	}

	if adminUsers != nil && len(adminUsers) > 0 && adminUsers[0].OpenID == "" {
		adminUser := adminUsers[0]
		adminUser.OpenID = openID
		err = us.adminUserModel.UpdateAdminUserByConditionWithTx(tx, adminUser, "id", "open_id")
		if err != nil {
			logger.Warn(userServiceLogTag, "bindAdminUser Update Failed|Uid:%v|Admin:%v|Err:%v",
				adminUser.ID, err)
			return err
		}
	}
	return nil
}

func (us *UserService) BindPhoneNumber(uid uint32, phoneNumber string) error {
	condition := " WHERE `id` = ?  "
	users, err := us.wxUserModel.GetWxUserByCondition(condition, uid)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetWxUserByCondition Failed|uid:%v|Err:%v", uid, err)
		return err
	}
	if len(users) == 0 {
		logger.Warn(userServiceLogTag, "User Not Exist|uid:%v|Err:%v", uid, err)
		return fmt.Errorf("user not exist|uid:%v", uid)
	}
	wxUser := users[0]

	if wxUser.PhoneNumber != "" {
		return fmt.Errorf("用户已经绑定过电话了，如需修改，请联系管理员")
	}

	_, err = us.CheckPhoneNumber(phoneNumber, false)
	if err != nil {
		return err
	}

	wxUser.PhoneNumber = phoneNumber
	tx, err := utils.Begin(us.sqlCli)
	if err != nil {
		logger.Warn(userServiceLogTag, "BindPhoneNumber Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	err = us.wxUserModel.UpdateWithTx(tx, wxUser, "id", "phone_number")
	if err != nil {
		logger.Warn(userServiceLogTag, "Update PhoneNumber Failed|Uid:%v|PhoneNumber:%v|Err:%v",
			uid, phoneNumber, err)
		return err
	}

	err = us.bindOrderUser(tx, wxUser.OpenID, wxUser.PhoneNumber, wxUser.ID)
	if err != nil {
		return err
	}
	err = us.bindAdminUser(tx, wxUser.OpenID, wxUser.PhoneNumber)
	if err != nil {
		return err
	}
	err = us.bindSupplier(tx, wxUser.OpenID, wxUser.PhoneNumber, uid)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) GetOrderUserList(openID, phoneNumber string, discountLevel, page, pageSize int32) ([]*model.OrderUser, uint32, error) {
	userList, err := us.orderUserModel.GetOrderUser(openID, phoneNumber, discountLevel, page, pageSize)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetOrderUser Failed|Err:%v", err)
		return nil, 0, err
	}

	userNumber, err := us.orderUserModel.GetOrderUserCount(openID, phoneNumber, discountLevel)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetOrderUserCount Failed|Err:%v", err)
		return nil, 0, err
	}
	return userList, uint32(userNumber), nil
}

func (us *UserService) AddOrderUser(userList []*model.OrderUser) error {
	err := us.orderUserModel.BatchInsert(userList)
	if err != nil {
		logger.Warn(userServiceLogTag, "AddOrderUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) BindOrderUser(id uint32, openID string) error {
	uid := uint32(0)
	if openID != "" {
		wxUser, err := us.wxUserModel.GetWxUserByOpenID(openID)
		if err != nil || wxUser == nil {
			logger.Warn(userServiceLogTag, "BindOrderUser GetWxUser Failed|Err:%v", err)
			return fmt.Errorf("用户不存在")
		}

		orderUserList, err := us.orderUserModel.GetOrderUserByCondition(" WHERE `open_id`=? ", wxUser.OpenID)
		if err != nil {
			logger.Warn(userServiceLogTag, "BindOrderUser GetOrderUser Failed|Err:%v", err)
			return err
		}
		if orderUserList != nil && len(orderUserList) > 0 {
			return fmt.Errorf("改微信已绑定，请解绑后重试")
		}
		uid = wxUser.ID
	}

	orderUser := &model.OrderUser{ID: id, OpenID: openID, Uid: uid}
	err := us.orderUserModel.UpdateOrderUserWithTx(nil, orderUser, "id", "open_id", "uid")
	if err != nil {
		logger.Warn(userServiceLogTag, "BindOrderUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) UpdateOrderUser(updateInfo *model.OrderUser) error {
	condition := " WHERE id = ? "
	orderUsers, err := us.orderUserModel.GetOrderUserByCondition(condition, updateInfo.ID)
	if err != nil {
		logger.Warn(userServiceLogTag, "UpdateOrderUser GetUser Failed|Err:%v", err)
		return err
	}
	if len(orderUsers) == 0 {
		return fmt.Errorf("订餐用户不存在")
	}

	updateTags := make([]string, 0)
	if updateInfo.OpenID != "" {
		wxUser, err := us.wxUserModel.GetWxUserByOpenID(updateInfo.OpenID)
		if err != nil {
			logger.Warn(userServiceLogTag, "GetWxUserByOpenID Failed|Err:%v", err)
			return err
		}
		updateInfo.Uid = wxUser.ID
		updateTags = append(updateTags, "open_id", "uid")
	}
	updateTags = append(updateTags, "discount_level")

	err = us.orderUserModel.UpdateOrderUserWithTx(nil, updateInfo, "id", updateTags...)
	if err != nil {
		logger.Warn(userServiceLogTag, "ModifyOrderUser Failed|Err:%v", err)
		return err
	}

	return nil
}

func (us *UserService) DeleteOrderUser(orderUseID uint32) error {
	err := us.orderUserModel.DeleteOrderUser(orderUseID)
	if err != nil {
		logger.Warn(userServiceLogTag, "DeleteOrderUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us UserService) GetAdminMap() (map[uint32]*model.AdminUser, error) {
	adminList, err := us.adminUserModel.GetAdminUserByCondition(" ORDER BY `id` ASC ")
	if err != nil {
		logger.Warn(userServiceLogTag, "GetAdminMap Failed|Err:%v", err)
		return nil, err
	}
	retMap := make(map[uint32]*model.AdminUser)
	for _, admin := range adminList {
		retMap[admin.ID] = admin
	}
	return retMap, nil
}

func (us UserService) GetAdminUserList(roleType uint8, page, pageSize int32) ([]*model.AdminUser, int32, error) {
	adminList, err := us.adminUserModel.GetAdminUserByCondition(" ORDER BY `id` ASC ")
	if err != nil {
		logger.Warn(userServiceLogTag, "GetAdminUserByCondition Failed|Err:%v", err)
		return nil, 0, err
	}
	if roleType == 0 {
		return adminList, int32(len(adminList)), nil
	}
	role := uint32(1 << roleType)
	retList := make([]*model.AdminUser, 0)
	for _, admin := range adminList {
		if (admin.Role & role) > 0 {
			retList = append(retList, admin)
		}
	}
	count := int32(len(retList))
	start := (page - 1) * pageSize
	if start >= count {
		retList = make([]*model.AdminUser, 0)
	} else {
		retList = retList[start:]
	}
	if pageSize < int32(len(retList)) {
		retList = retList[0:pageSize]
	}
	return retList, count, nil
}

func (us *UserService) AddAdminUser(user *model.AdminUser) error {
	if user.OpenID != "" {
		wxUser, err := us.wxUserModel.GetWxUserByOpenID(user.OpenID)
		if err != nil {
			logger.Warn(userServiceLogTag, "GetWxUserByOpenID Failed|Err:%v", err)
			return err
		}
		if wxUser == nil {
			logger.Warn(userServiceLogTag, "WxUser NotExist|OpenID:%v", user.OpenID)
			return fmt.Errorf("用户不存在，请确认OpenID是否正确")
		}
	} else if user.PhoneNumber != "" {
		wxUser, _ := us.CheckPhoneNumber(user.PhoneNumber, true)
		if wxUser != nil {
			user.OpenID = wxUser.OpenID
		}
	}
	err := us.adminUserModel.Insert(user)
	if err != nil {
		logger.Warn(userServiceLogTag, "Insert AdminUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) UpdateAdminUser(user *model.AdminUser) error {
	err := us.adminUserModel.UpdateAdminUserInfo(user, "id")
	if err != nil {
		logger.Warn(userServiceLogTag, "Update AdminUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) DeleteAdminUser(id uint32) error {
	err := us.adminUserModel.DeleteByID(id)
	if err != nil {
		logger.Warn(userServiceLogTag, "DeleteAdminUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) BindAdminUser(userID uint32, openID string) error {
	if openID != "" {
		wxUser, err := us.wxUserModel.GetWxUserByOpenID(openID)
		if err != nil {
			logger.Warn(userServiceLogTag, "GetWxUserByOpenID Failed|Err:%v", err)
			return err
		}
		if wxUser == nil {
			logger.Warn(userServiceLogTag, "WxUser NotExist|OpenID:%v", openID)
			return fmt.Errorf("要绑定的用户不存在，请确认OpenID是否正确")
		}
	}

	adminUser := &model.AdminUser{ID: userID, OpenID: openID}
	err := us.adminUserModel.UpdateAdminUserByCondition(adminUser, "id", "open_id")
	if err != nil {
		logger.Warn(userServiceLogTag, "BindAdminUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) GetRouterTypeList() ([]*model.RouterType, error) {
	typeList, err := us.routerTypeModel.GetRouterTypes()
	if err != nil {
		logger.Warn(userServiceLogTag, "GetGoodsTypeList Failed|Err:%v", err)
		return nil, err
	}
	return typeList, nil
}

func (us *UserService) GetRouterList(routerType uint32) ([]*model.RouterDetail, error) {
	routerList, err := us.routerDetailModel.GetRouterDetail(routerType)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetRouterList Failed|Err:%v", err)
		return nil, err
	}
	return routerList, nil
}

func (us *UserService) AddRouterType(routerType *model.RouterType) error {
	err := us.routerTypeModel.Insert(routerType)
	if err != nil {
		logger.Warn(userServiceLogTag, "AddRouterType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) UpdateRouterType(routerType *model.RouterType) error {
	err := us.routerTypeModel.UpdateRouterType(routerType)
	if err != nil {
		logger.Warn(userServiceLogTag, "UpdateRouterType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) DelRouterType(routerTypeID uint32) error {
	list, err := us.routerDetailModel.GetRouterDetail(routerTypeID)
	if err != nil {
		logger.Warn(userServiceLogTag, "DelRouterType GetRouterDetail Failed|Err:%v", err)
		return err
	}
	if len(list) > 0 {
		return fmt.Errorf("该路由类型下还有路由，无法删除")
	}
	err = us.routerTypeModel.DeleteRouterType(routerTypeID)
	if err != nil {
		logger.Warn(userServiceLogTag, "DeleteRouterType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) AddRouter(router *model.RouterDetail) error {
	err := us.routerDetailModel.Insert(router)
	if err != nil {
		logger.Warn(userServiceLogTag, "AddRouter Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) UpdateRouter(router *model.RouterDetail) error {
	err := us.routerDetailModel.UpdateDetail(router)
	if err != nil {
		logger.Warn(userServiceLogTag, "UpdateRouter Failed|Err:%v", err)
		return err
	}
	return nil
}

func (us *UserService) DeleteRouter(routerID uint32) error {
	err := us.routerDetailModel.DeleteByID(routerID)
	if err != nil {
		logger.Warn(userServiceLogTag, "DeleteRouter Failed|Err:%v", err)
		return err
	}
	return nil
}
