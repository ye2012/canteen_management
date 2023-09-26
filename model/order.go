package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	orderTable = "order"

	orderLogTag = "OrderModel"
)

type OrderDao struct {
	ID uint32 `json:"id"`
	//OrderID      string          `json:"order_id"`
	PayOrderID   uint32    `json:"pay_order_id"`
	MealType     uint8     `json:"meal_type"`
	OrderDate    time.Time `json:"order_date"`
	Uid          uint32    `json:"uid"`
	Address      string    `json:"address"`
	TotalAmount  float64   `json:"total_amount"`
	PayAmount    float64   `json:"pay_amount"`
	DiscountType uint8     `json:"discount_type"`
	Status       uint8     `json:"status"`
	CreateAt     time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}

type OrderModel struct {
	sqlCli *sql.DB
}

func NewOrderModel(sqlCli *sql.DB) *OrderModel {
	return &OrderModel{
		sqlCli: sqlCli,
	}
}

func (om *OrderModel) InsertWithTx(tx *sql.Tx, dao *OrderDao) (err error) {
	id := int64(0)
	if tx != nil {
		id, err = utils.SqlInsert(tx, orderTable, dao, "id", "created_at", "updated_at")
	} else {
		id, err = utils.SqlInsert(om.sqlCli, orderTable, dao, "id", "created_at", "updated_at")
	}

	if err != nil {
		logger.Warn(orderLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (om *OrderModel) UpdateOrderInfoByID(tx *sql.Tx, order *OrderDao, updateTags ...string) (err error) {
	if tx != nil {
		err = utils.SqlUpdateWithUpdateTags(tx, orderTable, order, "id", updateTags...)
	} else {
		err = utils.SqlUpdateWithUpdateTags(om.sqlCli, orderTable, order, "id", updateTags...)
	}
	if err != nil {
		logger.Warn(orderLogTag, "UpdateOrderInfoByID Failed|Err:%v", err)
		return err
	}
	return nil
}

func (om *OrderModel) GetOrderListByPayOrder(payOrderList []uint32) ([]*OrderDao, error) {
	if len(payOrderList) == 0 {
		return make([]*OrderDao, 0), nil
	}
	payOrderStr := ""
	for _, payOrderID := range payOrderList {
		payOrderStr += fmt.Sprintf(",%v", payOrderID)
	}
	condition := fmt.Sprintf(" WHERE `pay_order_id` in (%v) ", payOrderStr[1:])
	retList, err := utils.SqlQuery(om.sqlCli, orderTable, &OrderDao{}, condition)
	if err != nil {
		logger.Warn(orderLogTag, "GetOrderListByPayOrder Failed|PayOrderID:%v|Err:%v", payOrderList, err)
		return nil, err
	}

	return retList.([]*OrderDao), nil
}

func (om *OrderModel) GenerateCondition(idList []uint32, uid uint32, status int8) (string, []interface{}) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if len(idList) > 0 {
		idStr := ""
		for _, id := range idList {
			idStr += fmt.Sprintf(",%v", id)
		}
		condition += fmt.Sprintf(" AND `id` in (%v) ", idStr)
	}
	if uid > 0 {
		condition += " AND `uid` = ? "
		params = append(params, uid)
	}
	if status != -1 {
		condition += " AND `status` = ? "
		params = append(params, status)
	}
	return condition, params
}

func (om *OrderModel) GetOrderList(idList []uint32, uid uint32, page, pageSize int32, status int8) ([]*OrderDao, error) {
	condition, params := om.GenerateCondition(idList, uid, status)
	condition += " ORDER BY `id` ASC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(om.sqlCli, orderTable, &OrderDao{}, condition, params...)
	if err != nil {
		logger.Warn(orderLogTag, "GetOrderList Failed|ID:%v|Uid:%v|Status:%v|Err:%v", idList, uid, status, err)
		return nil, err
	}

	return retList.([]*OrderDao), nil
}

func (om *OrderModel) GetOrderListCount(idList []uint32, uid uint32, status int8) (int32, error) {
	condition, params := om.GenerateCondition(idList, uid, status)

	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM `%v` %v", orderTable, condition)
	row := om.sqlCli.QueryRow(sqlStr, params...)
	var count int32 = 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (om *OrderModel) GetOrder(id uint32) (*OrderDao, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	retInfo := &OrderDao{}
	err := utils.SqlQueryRow(om.sqlCli, orderTable, retInfo, condition, params...)
	if err != nil {
		logger.Warn(orderLogTag, "GetOrder Failed|ID:%v|Err:%v", id, err)
		return nil, err
	}

	return retInfo, nil
}
