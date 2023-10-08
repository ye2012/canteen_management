package model

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	purchaseDetailTable  = "purchase_detail"
	purchaseDetailLogTag = "PurchaseDetailModel"
)

var (
	purchaseDetailUpdateTags = []string{"receive_amount"}
)

type PurchaseDetail struct {
	ID            uint32  `json:"id"`
	PurchaseID    uint32  `json:"purchase_id"`
	GoodsID       uint32  `json:"goods_id"`
	GoodsType     uint32  `json:"goods_type"`
	ExpectAmount  float64 `json:"expect_amount"`
	ReceiveAmount float64 `json:"receive_amount"`
	Price         float64 `json:"price"`
}

type PurchaseDetailModel struct {
	sqlCli *sql.DB
}

func NewPurchaseDetailModelWithDB(sqlCli *sql.DB) *PurchaseDetailModel {
	return &PurchaseDetailModel{
		sqlCli: sqlCli,
	}
}

func (pdm *PurchaseDetailModel) BatchInsertWithTx(tx *sql.Tx, goodsList []*PurchaseDetail) (err error) {
	if tx != nil {
		err = utils.SqlInsertBatch(tx, purchaseDetailTable, goodsList, "id")
	} else {
		err = utils.SqlInsertBatch(pdm.sqlCli, purchaseDetailTable, goodsList, "id")
	}
	if err != nil {
		logger.Warn(purchaseDetailLogTag, "Insert Failed|GoodsList:%+v|Err:%v", goodsList, err)
		return err
	}
	return nil
}

func (pdm *PurchaseDetailModel) BatchInsert(goodsList []*PurchaseDetail) error {
	return pdm.BatchInsertWithTx(nil, goodsList)
}

func (pdm *PurchaseDetailModel) GetDetail(orderID uint32) ([]*PurchaseDetail, error) {
	condition := " WHERE `purchase_order_id` = ?  "
	retList, err := utils.SqlQuery(pdm.sqlCli, purchaseDetailTable, &PurchaseDetail{}, condition, orderID)
	if err != nil {
		logger.Warn(purchaseDetailLogTag, "GetDetail Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*PurchaseDetail), nil
}

func (pdm *PurchaseDetailModel) BatchUpdateDetailWithTx(tx *sql.Tx, detailList []*PurchaseDetail) (err error) {
	daoList := make([]interface{}, 0)
	for _, detail := range detailList {
		daoList = append(daoList, detail)
	}
	if tx != nil {
		err = utils.SqlBatchUpdateTag(tx, purchaseDetailTable, daoList, "id", "receive_amount")
	} else {
		err = utils.SqlBatchUpdateTag(pdm.sqlCli, purchaseDetailTable, daoList, "id", "receive_amount")
	}
	if err != nil {
		logger.Warn(purchaseDetailLogTag, "BatchUpdateDetail Failed|Err:%v", err)
		return err
	}
	return nil
}

func (pdm *PurchaseDetailModel) BatchUpdateDetail(detailList []*PurchaseDetail) error {
	return pdm.BatchUpdateDetailWithTx(nil, detailList)
}

func (pdm *PurchaseDetailModel) GetPurchaseDetailByOrderList(purchaseIDList []uint32, goodsType uint32) ([]*PurchaseDetail, error) {
	if len(purchaseIDList) == 0 {
		return nil, fmt.Errorf("purchase len zero")
	}
	purchaseIDStr := ""
	for _, purchaseID := range purchaseIDList {
		purchaseIDStr += fmt.Sprintf(",%v", purchaseID)
	}
	condition := fmt.Sprintf(" WHERE `purchase_id` in (%v) ", purchaseIDStr[1:])
	if goodsType != 0 {
		condition += fmt.Sprintf(" AND `goods_type` = %v ", goodsType)
	}
	retList, err := utils.SqlQuery(pdm.sqlCli, purchaseDetailTable, &PurchaseDetail{}, condition)
	if err != nil {
		logger.Warn(purchaseDetailLogTag, "GetPurchaseDetailByOrderList Failed|Condition:%v|purchaseID:%#v|Err:%v",
			condition, purchaseIDList, err)
		return nil, err
	}

	return retList.([]*PurchaseDetail), nil
}
