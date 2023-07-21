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
	ExpectAmount    float64 `json:"expect_amount"`
	ReceiveAmount   float64 `json:"receive_amount"`
	Discount        float64 `json:"discount"`
	DealPrice       float64 `json:"deal_price"`
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

func (pdm *PurchaseDetailModel) UpdateDetail(goods []*PurchaseDetail) error {
	err := utils.SqlUpdateWithUpdateTags(pdm.sqlCli, purchaseDetailTable, goods, "id", purchaseDetailUpdateTags...)
	if err != nil {
		logger.Warn(purchaseDetailLogTag, "UpdateDetail Failed|Err:%v", err)
		return err
	}
	return nil
}
