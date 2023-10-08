package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	purchaseOrderTable = "purchase_order"

	purchaseOrderLogTag = "PurchaseOrderModel"
)

var (
	purchaseOrderUpdateTags = []string{"status"}
)

type PurchaseOrder struct {
	ID          uint32    `json:"id"`
	Supplier    uint32    `json:"supplier"`
	TotalAmount float64   `json:"total_amount"`
	PayAmount   float64   `json:"pay_amount"`
	Creator     uint32    `json:"creator"`
	Status      uint8     `json:"status"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

type PurchaseOrderModel struct {
	sqlCli *sql.DB
}

func NewPurchaseOrderModelWithDB(sqlCli *sql.DB) *PurchaseOrderModel {
	return &PurchaseOrderModel{
		sqlCli: sqlCli,
	}
}

func (pom *PurchaseOrderModel) Insert(dao *PurchaseOrder) error {
	return pom.InsertWithTx(nil, dao)
}

func (pom *PurchaseOrderModel) InsertWithTx(tx *sql.Tx, dao *PurchaseOrder) error {
	id, err := int64(0), error(nil)
	if tx != nil {
		id, err = utils.SqlInsert(tx, purchaseOrderTable, dao, "id", "created_at", "updated_at")
	} else {
		id, err = utils.SqlInsert(pom.sqlCli, purchaseOrderTable, dao, "id", "created_at", "updated_at")
	}
	if err != nil {
		logger.Warn(purchaseOrderLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (pom *PurchaseOrderModel) GenerateCondition(id uint32, status int8, supplier, creator uint32, startTime,
	endTime int64) (string, []interface{}) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	if supplier > 0 {
		condition += " AND `supplier` = ? "
		params = append(params, supplier)
	}
	if creator > 0 {
		condition += " AND `creator` = ? "
		params = append(params, creator)
	}
	if startTime > 0 {
		condition += " AND `created_at` >= ? "
		params = append(params, time.Unix(startTime, 0))
	}
	if endTime > startTime {
		condition += " AND `created_at` <= ? "
		params = append(params, time.Unix(endTime, 0))
	}
	if status != enum.PurchaseOrderStatusAll {
		condition += " AND `status` = ? "
		params = append(params, status)
	}
	return condition, params
}

func (pom *PurchaseOrderModel) GetPurchaseOrder(id uint32) (*PurchaseOrder, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	retInfo := &PurchaseOrder{}
	err := utils.SqlQueryRow(pom.sqlCli, purchaseOrderTable, retInfo, condition, params...)
	if err != nil {
		logger.Warn(orderLogTag, "GetOrder Failed|ID:%v|Err:%v", id, err)
		return nil, err
	}

	return retInfo, nil
}

func (pom *PurchaseOrderModel) GetPurchaseOrderList(id uint32, status int8, supplier, creator uint32, startTime,
	endTime int64, page, pageSize int32) ([]*PurchaseOrder, error) {
	condition, params := pom.GenerateCondition(id, status, supplier, creator, startTime, endTime)
	condition += " ORDER BY `id` ASC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(pom.sqlCli, purchaseOrderTable, &PurchaseOrder{}, condition, params...)
	if err != nil {
		logger.Warn(purchaseOrderLogTag, "GetPurchaseOrders Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*PurchaseOrder), nil
}

func (pom *PurchaseOrderModel) GetPurchaseOrderCount(id uint32, status int8, supplier, creator uint32, startTime,
	endTime int64) (int32, error) {
	condition, params := pom.GenerateCondition(id, status, supplier, creator, startTime, endTime)
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM `%v` %v", purchaseOrderTable, condition)
	row := pom.sqlCli.QueryRow(sqlStr, params...)
	var count int32 = 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pom *PurchaseOrderModel) UpdatePurchaseWithTx(tx *sql.Tx, dao *PurchaseOrder, updateTags ...string) (err error) {
	if tx != nil {
		err = utils.SqlUpdateWithUpdateTags(tx, purchaseOrderTable, dao, "id", updateTags...)
	} else {
		err = utils.SqlUpdateWithUpdateTags(pom.sqlCli, purchaseOrderTable, dao, "id", updateTags...)
	}
	if err != nil {
		logger.Warn(purchaseOrderLogTag, "UpdatePurchase Failed|Err:%v", err)
		return err
	}
	return nil
}

func (pom *PurchaseOrderModel) UpdatePurchase(dao *PurchaseOrder, updateTags ...string) error {
	return pom.UpdatePurchaseWithTx(nil, dao, updateTags...)
}

func (pom *PurchaseOrderModel) UpdatePurchaseStatus(dao *PurchaseOrder) error {
	err := utils.SqlUpdateWithUpdateTags(pom.sqlCli, purchaseOrderTable, dao,
		"id", "status")
	if err != nil {
		logger.Warn(purchaseOrderLogTag, "UpdatePurchaseStatus Failed|Err:%v", err)
		return err
	}
	return nil
}
