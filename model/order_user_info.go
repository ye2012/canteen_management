package model

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	orderUserTable = "order_user"

	orderUserLogTag = "OrderUser"
)

type OrderUser struct {
	ID            uint32 `json:"id"`
	OpenID        string `json:"open_id"`
	Uid           uint32 `json:"uid"`
	PhoneNumber   string `json:"phone_number"`
	DiscountLevel uint8  `json:"discount_level"`
}

var (
	orderUserUpdateTag = []string{"open_id", "phone_number", "discount_level"}
)

type OrderUserModel struct {
	sqlCli *sql.DB
}

func NewOrderUserModel(sqlCli *sql.DB) *OrderUserModel {
	return &OrderUserModel{
		sqlCli: sqlCli,
	}
}

func (oum *OrderUserModel) BatchInsert(userList []*OrderUser) error {
	err := utils.SqlInsertBatch(oum.sqlCli, orderUserTable, userList, "id")
	if err != nil {
		logger.Warn(orderUserLogTag, "BatchInsert Failed|UserList:%v|Err:%v", userList, err)
		return err
	}
	return nil
}

func (oum *OrderUserModel) GetOrderUserByCondition(condition string, params ...interface{}) ([]*OrderUser, error) {
	retList, err := utils.SqlQuery(oum.sqlCli, orderUserTable, &OrderUser{}, condition, params...)
	if err != nil {
		logger.Warn(orderUserLogTag, "GetOrderUserByCondition Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*OrderUser), nil
}

func (oum *OrderUserModel) GetOrderUser(phoneNumber string, discountLevel, page, pageSize int32) ([]*OrderUser, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if phoneNumber != "" {
		condition += " AND `phone_number` = ? "
		params = append(params, phoneNumber)
	}
	if discountLevel != 0 {
		condition += " AND `discount_level` = ? "
		params = append(params, discountLevel)
	}
	condition += " ORDER BY `id` ASC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(oum.sqlCli, orderUserTable, &OrderUser{}, condition, params...)
	if err != nil {
		logger.Warn(orderUserLogTag, "GetOrderUser Failed|phoneNumber:%v|Err:%v", phoneNumber, err)
		return nil, err
	}

	return retList.([]*OrderUser), nil
}

func (oum *OrderUserModel) GetOrderUserCount(phoneNumber string, discountLevel int32) (int32, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if phoneNumber != "" {
		condition += " AND `phone_number` = ? "
		params = append(params, phoneNumber)
	}
	if discountLevel != 0 {
		condition += " AND `discount_level` = ? "
		params = append(params, discountLevel)
	}
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM `%v` %v", orderUserTable, condition)
	row := oum.sqlCli.QueryRow(sqlStr, params...)
	var count int32
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (oum *OrderUserModel) UpdateOrderUser(userInfo *OrderUser) error {
	return oum.UpdateOrderUserWithTx(nil, userInfo, "id", orderUserUpdateTag...)
}

func (oum *OrderUserModel) UpdateOrderUserWithTx(tx *sql.Tx, userInfo *OrderUser, conditionTag string, params ...string) error {
	if tx == nil {
		err := utils.SqlUpdateWithUpdateTags(oum.sqlCli, orderUserTable, userInfo, conditionTag, params...)
		if err != nil {
			logger.Warn(orderUserLogTag, "SqlUpdateWithUpdateTags Failed|Err:%v", err)
			return err
		}
	} else {
		err := utils.SqlUpdateWithUpdateTags(tx, orderUserTable, userInfo, conditionTag, params...)
		if err != nil {
			logger.Warn(orderUserLogTag, "SqlUpdateWithUpdateTags Failed|Err:%v", err)
			return err
		}
	}
	return nil
}
