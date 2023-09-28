package model

import (
	"database/sql"
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
	id, err := utils.SqlInsert(pom.sqlCli, purchaseOrderTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(purchaseOrderLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (pom *PurchaseOrderModel) GetPurchaseOrder(id uint32, status int8, supplier uint32, startTime,
	endTime int64) ([]*PurchaseOrder, error) {
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
	retList, err := utils.SqlQuery(pom.sqlCli, purchaseOrderTable, &PurchaseOrder{}, condition, params...)
	if err != nil {
		logger.Warn(purchaseOrderLogTag, "GetPurchaseOrders Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*PurchaseOrder), nil
}

func (pom *PurchaseOrderModel) UpdatePurchase(dao *PurchaseOrder, updateTags ...string) error {
	err := utils.SqlUpdateWithUpdateTags(pom.sqlCli, purchaseOrderTable, dao,
		"id", updateTags...)
	if err != nil {
		logger.Warn(purchaseOrderLogTag, "UpdatePurchase Failed|Err:%v", err)
		return err
	}
	return nil
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
