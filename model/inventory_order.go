package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	inventoryOrderTable  = "inventory_order"
	inventoryModelLogTag = "InventoryModel"
)

type InventoryOrder struct {
	ID       uint32    `json:"id"`
	Creator  uint32    `json:"creator"`
	Partner  uint32    `json:"partner"`
	Status   int8      `json:"status"`
	CreateAt time.Time `json:"created_at"`
	FinishAt time.Time `json:"finish_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type InventoryOrderModel struct {
	sqlCli *sql.DB
}

func NewInventoryOrderModel(sqlCli *sql.DB) *InventoryOrderModel {
	return &InventoryOrderModel{
		sqlCli: sqlCli,
	}
}

func (iom *InventoryOrderModel) Insert(dao *InventoryOrder) error {
	return iom.InsertWithTx(nil, dao)
}

func (iom *InventoryOrderModel) InsertWithTx(tx *sql.Tx, dao *InventoryOrder) error {
	id, err := int64(0), error(nil)
	if tx != nil {
		id, err = utils.SqlInsert(tx, inventoryOrderTable, dao, "id", "created_at", "updated_at", "finish_at")
	} else {
		id, err = utils.SqlInsert(iom.sqlCli, inventoryOrderTable, dao, "id", "created_at", "updated_at", "finish_at")
	}
	if err != nil {
		logger.Warn(inventoryModelLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (iom *InventoryOrderModel) GenerateCondition(id, creator uint32, status int8, startTime, endTime int64) (string, []interface{}) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	if creator > 0 {
		condition += " AND `creator` = ? "
		params = append(params, creator)
	}
	if status != -1 {
		condition += " AND `status` = ? "
		params = append(params, status)
	}
	if startTime > 0 {
		condition += " AND `created_at` >= ? "
		params = append(params, time.Unix(startTime, 0))
	}
	if endTime > startTime {
		condition += " AND `created_at` <= ? "
		params = append(params, time.Unix(endTime, 0))
	}
	return condition, params
}

func (iom *InventoryOrderModel) GetInventoryOrder(id uint32) (*InventoryOrder, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	retInfo := &InventoryOrder{}
	err := utils.SqlQueryRow(iom.sqlCli, inventoryOrderTable, retInfo, condition, params...)
	if err != nil {
		logger.Warn(inventoryModelLogTag, "GetInventoryOrder Failed|ID:%v|Err:%v", id, err)
		return nil, err
	}

	return retInfo, nil
}

func (iom *InventoryOrderModel) GetInventoryOrderList(id, creator uint32, status int8, startTime, endTime int64,
	page, pageSize int32) ([]*InventoryOrder, error) {
	condition, params := iom.GenerateCondition(id, creator, status, startTime, endTime)
	condition += " ORDER BY `id` DESC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(iom.sqlCli, inventoryOrderTable, &InventoryOrder{}, condition, params...)
	if err != nil {
		logger.Warn(inventoryModelLogTag, "GetInventoryOrders Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*InventoryOrder), nil
}

func (iom *InventoryOrderModel) GetInventoryOrderCount(id, creator uint32, status int8, startTime,
	endTime int64) (int32, error) {
	condition, params := iom.GenerateCondition(id, creator, status, startTime, endTime)
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM `%v` %v", inventoryOrderTable, condition)
	row := iom.sqlCli.QueryRow(sqlStr, params...)
	var count int32 = 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (iom *InventoryOrderModel) UpdateInventoryOrderByCondition(order *InventoryOrder, conditionTag string,
	updateTags ...string) error {
	return iom.UpdateInventoryOrderByConditionWithTx(nil, order, conditionTag, updateTags...)
}

func (iom *InventoryOrderModel) UpdateInventoryOrderByConditionWithTx(tx *sql.Tx, order *InventoryOrder, conditionTag string,
	updateTags ...string) (err error) {
	if tx != nil {
		err = utils.SqlUpdateWithUpdateTags(tx, inventoryOrderTable, order, conditionTag, updateTags...)
	} else {
		err = utils.SqlUpdateWithUpdateTags(iom.sqlCli, inventoryOrderTable, order, conditionTag, updateTags...)
	}
	if err != nil {
		logger.Warn(inventoryModelLogTag, "UpdateInventoryOrderByCondition Failed|Err:%v", err)
		return err
	}
	return nil
}
