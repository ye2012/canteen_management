package model

import (
	"database/sql"

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
	Id              int64   `json:"id"`
	PurchaseOrderID uint32  `json:"purchase_order_id"`
	GoodsID         uint32  `json:"goods_id"`
	GoodsType       uint32  `json:"goods_type"`
	ExpectAmount    float64 `json:"expect_amount"`
	ReceiveAmount   float64 `json:"receive_amount"`
	Price           float64 `json:"price"`
}

type PurchaseDetailModel struct {
	sqlCli *sql.DB
}

func NewPurchaseDetailModelWithDB(sqlCli *sql.DB) *PurchaseDetailModel {
	return &PurchaseDetailModel{
		sqlCli: sqlCli,
	}
}

func (pdm *PurchaseDetailModel) BatchInsert(goodsList []*PurchaseDetail) error {
	err := utils.SqlInsertBatch(pdm.sqlCli, purchaseDetailTable, goodsList, "id")
	if err != nil {
		logger.Warn(purchaseDetailLogTag, "Insert Failed|GoodsList:%+v|Err:%v", goodsList, err)
		return err
	}
	return nil
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

func (pdm *PurchaseDetailModel) BatchUpdateDetail(detailList []*PurchaseDetail) error {
	daoList := make([]interface{}, 0)
	for _, detail := range detailList {
		daoList = append(daoList, detail)
	}
	err := utils.SqlBatchUpdateTag(pdm.sqlCli, purchaseDetailTable, daoList, "id", "receive_amount")
	if err != nil {
		logger.Warn(purchaseDetailLogTag, "BatchUpdateDetail Failed|Err:%v", err)
		return err
	}
	return nil
}
