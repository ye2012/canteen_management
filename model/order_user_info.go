package model

import (
	"database/sql"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	orderUserTable = "order_user"

	orderUserLogTag = "OrderUser"
)

type OrderUserInfo struct {
	ID            uint32 `json:"id"`
	PhoneNumber   string `json:"phone_number"`
	DiscountLevel int32  `json:"discount_level"`
}

var (
	orderUserUpdateTag = []string{"phone_number", "discount_level"}
)

type OrderUserModel struct {
	sqlCli *sql.DB
}

func NewOrderUserModel(sqlCli *sql.DB) *OrderUserModel {
	return &OrderUserModel{
		sqlCli: sqlCli,
	}
}

func (oum *OrderUserModel) BatchInsert(userList []*OrderUserInfo) error {
	err := utils.SqlInsertBatch(oum.sqlCli, orderUserTable, userList, "id")
	if err != nil {
		logger.Warn(orderUserLogTag, "BatchInsert Failed|UserList:%v|Err:%v", userList, err)
		return err
	}
	return nil
}

func (oum *OrderUserModel) GetOrderUserByCondition(condition string, params ...interface{}) ([]*OrderUserInfo, error) {
	retList, err := utils.SqlQuery(oum.sqlCli, orderUserTable, &OrderUserInfo{}, condition, params...)
	if err != nil {
		logger.Warn(orderUserLogTag, "GetOrderUserByCondition Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*OrderUserInfo), nil
}

func (oum *OrderUserModel) GetOrderUser(phoneNumber string, page, pageSize uint32) ([]*OrderUserInfo, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if phoneNumber != "" {
		condition += " `phone_number` = ? "
		params = append(params, phoneNumber)
	}
	condition += " ORDER BY `id` ASC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(oum.sqlCli, orderUserTable, &OrderUserInfo{}, condition, params...)
	if err != nil {
		logger.Warn(orderUserLogTag, "GetOrderUser Failed|phoneNumber:%v|Err:%v", phoneNumber, err)
		return nil, err
	}

	return retList.([]*OrderUserInfo), nil
}

func (oum *OrderUserModel) UpdateOrderUser(userInfo *OrderUserInfo) error {
	err := utils.SqlUpdateWithUpdateTags(oum.sqlCli, orderUserTable, userInfo, "id", orderUserUpdateTag...)
	if err != nil {
		logger.Warn(goodsLogTag, "UpdateOrderUser Failed|Err:%v", err)
		return err
	}
	return nil
}
