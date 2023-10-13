package model

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	inventoryDetailTable  = "inventory_detail"
	inventoryDetailLogTag = "InventoryDetailModel"
)

type InventoryDetail struct {
	ID           uint32  `json:"id"`
	InventoryID  uint32  `json:"inventory_id"`
	GoodsID      uint32  `json:"goods_id"`
	GoodsType    uint32  `json:"goods_type"`
	ExpectNumber float64 `json:"expect_number"`
	RealNumber   float64 `json:"real_number"`
	Tag          string  `json:"tag"`
	Status       int8    `json:"status"`
}

type InventoryDetailModel struct {
	sqlCli *sql.DB
}

func NewInventoryDetailModelWithDB(sqlCli *sql.DB) *InventoryDetailModel {
	return &InventoryDetailModel{
		sqlCli: sqlCli,
	}
}

func (odm *InventoryDetailModel) BatchInsertWithTx(tx *sql.Tx, goodsList []*InventoryDetail) (err error) {
	if tx != nil {
		err = utils.SqlInsertBatch(tx, inventoryDetailTable, goodsList, "id")
	} else {
		err = utils.SqlInsertBatch(odm.sqlCli, inventoryDetailTable, goodsList, "id")
	}
	if err != nil {
		logger.Warn(inventoryDetailLogTag, "Insert Failed|GoodsList:%+v|Err:%v", goodsList, err)
		return err
	}
	return nil
}

func (odm *InventoryDetailModel) BatchInsert(goodsList []*InventoryDetail) error {
	return odm.BatchInsertWithTx(nil, goodsList)
}

func (odm *InventoryDetailModel) GetDetail(orderID uint32, status int8) ([]*InventoryDetail, error) {
	condition := " WHERE `inventory_id` = ?  "
	params := []interface{}{orderID}
	if status != -1 {
		condition += " AND `status` = ? "
		params = append(params, status)
	}
	retList, err := utils.SqlQuery(odm.sqlCli, inventoryDetailTable, &InventoryDetail{}, condition, params...)
	if err != nil {
		logger.Warn(inventoryDetailLogTag, "GetDetail Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*InventoryDetail), nil
}

func (odm *InventoryDetailModel) GetInventoryDetailByOrderList(inventoryIDList []uint32, goodsType uint32) ([]*InventoryDetail, error) {
	if len(inventoryIDList) == 0 {
		return nil, fmt.Errorf("inventory len zero")
	}
	inventoryIDStr := ""
	for _, inventoryID := range inventoryIDList {
		inventoryIDStr += fmt.Sprintf(",%v", inventoryID)
	}
	condition := fmt.Sprintf(" WHERE `inventory_id` in (%v) ", inventoryIDStr[1:])
	if goodsType != 0 {
		condition += fmt.Sprintf(" AND `goods_type` = %v ", goodsType)
	}
	retList, err := utils.SqlQuery(odm.sqlCli, inventoryDetailTable, &InventoryDetail{}, condition)
	if err != nil {
		logger.Warn(inventoryDetailLogTag, "GetInventoryDetailByOrderList Failed|Condition:%v|inventoryID:%#v|Err:%v",
			condition, inventoryIDList, err)
		return nil, err
	}

	return retList.([]*InventoryDetail), nil
}

func (odm *InventoryDetailModel) UpdateDetailByCondition(detail *InventoryDetail, conditionTag string,
	updateTags ...string) error {
	err := utils.SqlUpdateWithUpdateTags(odm.sqlCli, inventoryDetailTable, detail, conditionTag, updateTags...)
	if err != nil {
		logger.Warn(inventoryDetailLogTag, "UpdateDetailByCondition Failed|Err:%v", err)
		return err
	}
	return nil
}
