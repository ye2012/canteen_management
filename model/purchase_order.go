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
	SignPicture string    `json:"sign_picture"`
	Status      uint8     `json:"status"`
	Creator     uint32    `json:"creators"`
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

func (pom *PurchaseOrderModel) GetPurchaseOrder(status int8, supplier uint32, startTime,
	endTime time.Time) ([]*PurchaseOrder, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if supplier > 0 {
		condition += " AND `supplier` = ? "
		params = append(params, supplier)
	}
	if startTime.Unix() > 0 {
		condition += " AND `created_at` >= ? "
		params = append(params, startTime)
	}
	if endTime.Unix() > startTime.Unix() {
		condition += " AND `created_at` <= ? "
		params = append(params, endTime)
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

func (pom *PurchaseOrderModel) UpdatePurchase(dao *PurchaseOrder) error {
	err := utils.SqlUpdateWithUpdateTags(pom.sqlCli, purchaseOrderTable, dao,
		"id", purchaseOrderUpdateTags...)
	if err != nil {
		logger.Warn(purchaseOrderLogTag, "UpdatePurchase Failed|Err:%v", err)
		return err
	}
	return nil
}
