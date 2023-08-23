package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	orderTable = "order"

	orderLogTag = "OrderModel"
)

var (
	orderUpdateTags = []string{"status"}
)

type OrderDao struct {
	ID uint32 `json:"id"`
	//OrderID      string          `json:"order_id"`
	PrepareID    string    `json:"prepare_id"`
	MealType     uint8     `json:"meal_type"`
	OrderDate    time.Time `json:"order_date"`
	Uid          uint32    `json:"uid"`
	UnionID      string    `json:"union_id"`
	Address      string    `json:"address"`
	PickUpMethod uint8     `json:"pick_up_method"`
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

func (om *OrderModel) GetOrderList(id, uid, page, pageSize uint32, status int8) ([]*OrderDao, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	if uid > 0 {
		condition += " AND `uid` = ? "
		params = append(params, uid)
	}
	if status != -1 {
		condition += " AND `status` = ? "
		params = append(params, status)
	}
	condition += " ORDER BY `id` ASC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(om.sqlCli, orderTable, &OrderDao{}, condition, params...)
	if err != nil {
		logger.Warn(orderLogTag, "GetOrderList Failed|ID:%v|Uid:%v|Status:%v|Err:%v", id, uid, status, err)
		return nil, err
	}

	return retList.([]*OrderDao), nil
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
