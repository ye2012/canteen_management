package service

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/utils"
)

const (
	userServiceLogTag = "UserService"
)

type UserService struct {
	sqlCli         *sql.DB
	adminUserModel *model.AdminUserModel
	wxUserModel    *model.WxUserModel
	orderUserModel *model.OrderUserModel
}

func NewUserService(sqlCli *sql.DB) *UserService {
	wxUserModel := model.NewWxUserModelWithDB(sqlCli)
	orderUserModel := model.NewOrderUserModel(sqlCli)
	adminUserModel := model.NewAdminUserModelWithDB(sqlCli)
	return &UserService{
		sqlCli:         sqlCli,
		wxUserModel:    wxUserModel,
		orderUserModel: orderUserModel,
		adminUserModel: adminUserModel,
	}
}

func (us *UserService) GetWxUserRole(openID string) uint32 {
	adminList, err := us.adminUserModel.GetAdminUserByCondition(" WHERE `open_id`=? ", openID)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetAdminUserByCondition Failed|OpenID:%v|Err:%v", openID, err)
		return 0
	}
	if len(adminList) == 0 {
		return 0
	}
	return adminList[0].Role
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

func (us *UserService) WxUserLogin(openID string) (*model.WxUser, error) {
	user, err := us.wxUserModel.GetWxUserByOpenID(openID)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetWxUserByOpenID Failed|OpenID:%v|Err:%v", openID, err)
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	wxUser := &model.WxUser{OpenID: openID}
	err = us.wxUserModel.Insert(wxUser)
	if err != nil {
		logger.Warn(userServiceLogTag, "Insert WxUser Failed|OpenID:%v|Err:%v", openID, err)
		return nil, err
	}
	return wxUser, nil
}

func (us *UserService) BindPhoneNumber(uid uint32, phoneNumber string) error {
	condition := " WHERE `id` = ? "
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

	orderUser, err := us.orderUserModel.GetOrderUser(phoneNumber, 0, 1, 10)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetOrderUser Failed|phone:%v|Err:%v", phoneNumber, err)
		return err
	}
	discountType := uint8(0)
	if len(orderUser) > 0 {
		discountType = uint8(orderUser[0].DiscountLevel)
	}

	wxUser.OrderDiscountType = discountType
	wxUser.PhoneNumber = phoneNumber
	tx, err := utils.Begin(us.sqlCli)
	if err != nil {
		logger.Warn(userServiceLogTag, "BindPhoneNumber Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	err = us.wxUserModel.UpdateWithTx(tx, wxUser, "id", "phone_number", "order_discount_type")
	if err != nil {
		logger.Warn(userServiceLogTag, "Update PhoneNumber Failed|Uid:%v|PhoneNumber:%v|Err:%v",
			uid, phoneNumber, err)
		return err
	}

	if len(orderUser) > 0 {
		orderUser[0].OpenID = wxUser.OpenID
		orderUser[0].Uid = wxUser.ID
		err = us.orderUserModel.UpdateOrderUserWithTx(tx, orderUser[0], "id", "open_id", "uid")
		if err != nil {
			logger.Warn(userServiceLogTag, "UpdateOrderUserWithTx Failed|Uid:%v|OrderUser:%v|Err:%v",
				uid, orderUser[0].ID, err)
			return err
		}
	}

	return nil
}

func (us *UserService) GetOrderUserList(phoneNumber string, discountLevel, page, pageSize int32) ([]*model.OrderUser, uint32, error) {
	userList, err := us.orderUserModel.GetOrderUser(phoneNumber, discountLevel, page, pageSize)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetOrderUser Failed|Err:%v", err)
		return nil, 0, err
	}

	userNumber, err := us.orderUserModel.GetOrderUserCount(phoneNumber, discountLevel)
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

	tx, err := utils.Begin(us.sqlCli)
	if err != nil {
		logger.Warn(userServiceLogTag, "UpdateOrderUser Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	wxUser := &model.WxUser{ID: orderUsers[0].Uid, OrderDiscountType: updateInfo.DiscountLevel}
	updateTags := make([]string, 0)
	if orderUsers[0].OpenID != updateInfo.OpenID {
		wxUser, err = us.wxUserModel.GetWxUserByOpenID(updateInfo.OpenID)
		if err != nil {
			logger.Warn(userServiceLogTag, "GetWxUserByOpenID Failed|Err:%v", err)
			return err
		}
		err = us.wxUserModel.UpdateWithTx(tx, wxUser, "id", "order_discount_type")
		if err != nil {
			logger.Warn(userServiceLogTag, "Update WxUser Discount Failed|Err:%v", err)
			return err
		}
		updateInfo.Uid = wxUser.ID
		updateTags = append(updateTags, "open_id", "uid")
	}
	if orderUsers[0].DiscountLevel != updateInfo.DiscountLevel {
		if orderUsers[0].Uid != 0 {
			err = us.wxUserModel.UpdateWithTx(tx, wxUser, "id", "order_discount_type")
			if err != nil {
				logger.Warn(userServiceLogTag, "Update WxUser Discount Failed|Err:%v", err)
				return err
			}
		}
		updateTags = append(updateTags, "discount_level")
	}

	if len(updateTags) == 0 {
		return nil
	}

	err = us.orderUserModel.UpdateOrderUserWithTx(tx, updateInfo, "id", updateTags...)
	if err != nil {
		logger.Warn(userServiceLogTag, "ModifyOrderUser Failed|Err:%v", err)
		return err
	}

	return nil
}

func (us *UserService) AddAdminUser(user *model.AdminUser) {
	us.adminUserModel.Insert(user)
}

func (us *UserService) UpdateAdminUser(user *model.AdminUser) {
	us.adminUserModel.UpdateAdminUserInfo(user, "id")
}

func (us *UserService) DeleteAdminUser(user *model.AdminUser) {
	us.adminUserModel.Insert(user)
}

func (us *UserService) BindAdminUser(userID uint32, openID string) {

}
