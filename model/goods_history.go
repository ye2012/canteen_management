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
	goodsHistoryTable = "goods_history"

	goodsHistoryLogTag = "GoodsHistory"
)

type GoodsHistory struct {
	ID             uint32    `json:"id"`
	GoodsID        uint32    `json:"goods_id"`
	ChangeQuantity float64   `json:"change_quantity"`
	BeforeQuantity float64   `json:"before_quantity"`
	AfterQuantity  float64   `json:"after_quantity"`
	ChangeType     uint32    `json:"change_type"`
	RefID          uint32    `json:"ref_id"`
	CreateAt       time.Time `json:"created_at"`
}

func GenerateGoodsHistory(goodsID uint32, preQuantity, changeQuantity float64, changeType, refID uint32) *GoodsHistory {
	return &GoodsHistory{
		GoodsID:        goodsID,
		ChangeQuantity: changeQuantity,
		BeforeQuantity: preQuantity,
		AfterQuantity:  preQuantity + changeQuantity,
		ChangeType:     changeType,
		RefID:          refID,
	}
}

func GenerateInitGoodsHistory(goods *Goods) *GoodsHistory {
	return GenerateGoodsHistory(goods.ID, 0, goods.Quantity, enum.GoodsInit, 0)
}

func GenerateInventoryGoodsHistory(goods *Goods, changeQuantity float64, inventoryID uint32) *GoodsHistory {
	return GenerateGoodsHistory(goods.ID, goods.Quantity, changeQuantity, enum.GoodsInventory, inventoryID)
}

func GeneratePurchaseGoodsHistory(goods *Goods, changeQuantity float64, purchaseID uint32) *GoodsHistory {
	return GenerateGoodsHistory(goods.ID, goods.Quantity, changeQuantity, enum.GoodsPurchase, purchaseID)
}

func GenerateOutboundGoodsHistory(goods *Goods, changeQuantity float64, outboundID uint32) *GoodsHistory {
	return GenerateGoodsHistory(goods.ID, goods.Quantity, changeQuantity, enum.GoodsOutbound, outboundID)
}

type GoodsHistoryModel struct {
	sqlCli *sql.DB
}

func NewGoodsHistoryModel(sqlCli *sql.DB) *GoodsHistoryModel {
	return &GoodsHistoryModel{
		sqlCli: sqlCli,
	}
}

func (ghm *GoodsHistoryModel) BatchInsert(tx *sql.Tx, goodsHistory []*GoodsHistory) (err error) {
	if tx != nil {
		err = utils.SqlInsertBatch(tx, goodsHistoryTable, goodsHistory, "id", "created_at")
	} else {
		err = utils.SqlInsertBatch(ghm.sqlCli, goodsHistoryTable, goodsHistory, "id", "created_at")
	}
	if err != nil {
		logger.Warn(goodsHistoryLogTag, "BatchInsert Failed|GoodsHistory:%+v|Err:%v", goodsHistory, err)
		return err
	}
	return nil
}

func (ghm *GoodsHistoryModel) GenerateCondition(goodsID, changeType uint32, startTime, endTime int64) (string, []interface{}) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if goodsID > 0 {
		condition += " AND `goods_id` = ? "
		params = append(params, goodsID)
	}
	if changeType > 0 {
		condition += " AND `change_type` = ? "
		params = append(params, changeType)
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

func (ghm *GoodsHistoryModel) GetGoodsHistory(goodsID, changeType uint32, startTime, endTime int64,
	page, pageSize int32) ([]*GoodsHistory, error) {
	condition, params := ghm.GenerateCondition(goodsID, changeType, startTime, endTime)
	condition += " ORDER BY `id` DESC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(ghm.sqlCli, goodsHistoryTable, &GoodsHistory{}, condition, params...)
	if err != nil {
		logger.Warn(goodsHistoryLogTag, "GetGoodsHistory Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*GoodsHistory), nil
}

func (ghm *GoodsHistoryModel) GetGoodsHistoryCount(goodsID, changeType uint32, startTime, endTime int64) (int32, error) {
	condition, params := ghm.GenerateCondition(goodsID, changeType, startTime, endTime)
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM `%v` %v", goodsHistoryTable, condition)
	row := ghm.sqlCli.QueryRow(sqlStr, params...)
	var count int32 = 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ghm *GoodsHistoryModel) GetGoodsHistoryByCondition(condition string, params ...interface{}) ([]*GoodsHistory, error) {
	retList, err := utils.SqlQuery(ghm.sqlCli, goodsHistoryTable, &GoodsHistory{}, condition, params...)
	if err != nil {
		logger.Warn(goodsHistoryLogTag, "GetGoodsHistoryByCondition Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*GoodsHistory), nil
}
