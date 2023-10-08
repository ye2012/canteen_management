package model

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	orderDetailTable = "order_detail"

	orderDetailLogTag = "OrderDetail"
)

type OrderDetail struct {
	ID       uint32  `json:"id"`
	OrderID  uint32  `json:"order_id"`
	DishID   uint32  `json:"dish_id"`
	DishType uint32  `json:"dish_type"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

type OrderDetailModel struct {
	sqlCli *sql.DB
}

func NewOrderDetailModel(sqlCli *sql.DB) *OrderDetailModel {
	return &OrderDetailModel{
		sqlCli: sqlCli,
	}
}

func (odm *OrderDetailModel) BatchInsert(tx *sql.Tx, orderDetail []*OrderDetail) (err error) {
	if tx != nil {
		err = utils.SqlInsertBatch(tx, orderDetailTable, orderDetail, "id")
	} else {
		err = utils.SqlInsertBatch(odm.sqlCli, orderDetailTable, orderDetail, "id")
	}
	if err != nil {
		logger.Warn(orderDetailLogTag, "BatchInsert Failed|OrderDetail:%+v|Err:%v", orderDetail, err)
		return err
	}
	return nil
}

func (odm *OrderDetailModel) GetOrderDetailByCondition(condition string, params ...interface{}) ([]*OrderDetail, error) {
	retList, err := utils.SqlQuery(odm.sqlCli, orderDetailTable, &OrderDetail{}, condition, params...)
	if err != nil {
		logger.Warn(orderDetailLogTag, "GetOrderDetail Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*OrderDetail), nil
}

func (odm *OrderDetailModel) GetOrderDetail(orderID uint32) ([]*OrderDetail, error) {
	condition := " WHERE `order_id` = ? "
	retList, err := utils.SqlQuery(odm.sqlCli, orderDetailTable, &OrderDetail{}, condition, orderID)
	if err != nil {
		logger.Warn(orderDetailLogTag, "GetOrderDetail Failed|Condition:%v|OrderID:%#v|Err:%v",
			condition, orderID, err)
		return nil, err
	}

	return retList.([]*OrderDetail), nil
}

func (odm *OrderDetailModel) GetOrderDetailByOrderList(orderIDList []uint32, dishType uint32) ([]*OrderDetail, error) {
	if len(orderIDList) == 0 {
		return nil, fmt.Errorf("order len zero")
	}
	orderIDStr := ""
	for _, orderID := range orderIDList {
		orderIDStr += fmt.Sprintf(",%v", orderID)
	}
	condition := fmt.Sprintf(" WHERE `order_id` in (%v) ", orderIDStr[1:])
	if dishType != 0 {
		condition += fmt.Sprintf(" AND `dish_type` = %v ", dishType)
	}
	retList, err := utils.SqlQuery(odm.sqlCli, orderDetailTable, &OrderDetail{}, condition)
	if err != nil {
		logger.Warn(orderDetailLogTag, "GetOrderDetail Failed|Condition:%v|OrderID:%#v|Err:%v",
			condition, orderIDList, err)
		return nil, err
	}

	return retList.([]*OrderDetail), nil
}
