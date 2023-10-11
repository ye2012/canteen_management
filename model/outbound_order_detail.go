package model

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	outboundDetailTable  = "outbound_detail"
	outboundDetailLogTag = "OutboundDetailModel"
)

type OutboundDetail struct {
	ID         uint32  `json:"id"`
	OutboundID uint32  `json:"outbound_id"`
	GoodsID    uint32  `json:"goods_id"`
	GoodsType  uint32  `json:"goods_type"`
	OutNumber  float64 `json:"out_number"`
	Price      float64 `json:"price"`
}

type OutboundDetailModel struct {
	sqlCli *sql.DB
}

func NewOutboundDetailModelWithDB(sqlCli *sql.DB) *OutboundDetailModel {
	return &OutboundDetailModel{
		sqlCli: sqlCli,
	}
}

func (odm *OutboundDetailModel) BatchInsertWithTx(tx *sql.Tx, goodsList []*OutboundDetail) (err error) {
	if tx != nil {
		err = utils.SqlInsertBatch(tx, outboundDetailTable, goodsList, "id")
	} else {
		err = utils.SqlInsertBatch(odm.sqlCli, outboundDetailTable, goodsList, "id")
	}
	if err != nil {
		logger.Warn(outboundDetailLogTag, "Insert Failed|GoodsList:%+v|Err:%v", goodsList, err)
		return err
	}
	return nil
}

func (odm *OutboundDetailModel) BatchInsert(goodsList []*OutboundDetail) error {
	return odm.BatchInsertWithTx(nil, goodsList)
}

func (odm *OutboundDetailModel) GetDetail(orderID uint32) ([]*OutboundDetail, error) {
	condition := " WHERE `outbound_id` = ?  "
	retList, err := utils.SqlQuery(odm.sqlCli, outboundDetailTable, &OutboundDetail{}, condition, orderID)
	if err != nil {
		logger.Warn(outboundDetailLogTag, "GetDetail Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*OutboundDetail), nil
}

func (odm *OutboundDetailModel) GetOutboundDetailByOrderList(outboundIDList []uint32, goodsType uint32) ([]*OutboundDetail, error) {
	if len(outboundIDList) == 0 {
		return nil, fmt.Errorf("outbound len zero")
	}
	outboundIDStr := ""
	for _, outboundID := range outboundIDList {
		outboundIDStr += fmt.Sprintf(",%v", outboundID)
	}
	condition := fmt.Sprintf(" WHERE `outbound_id` in (%v) ", outboundIDStr[1:])
	if goodsType != 0 {
		condition += fmt.Sprintf(" AND `goods_type` = %v ", goodsType)
	}
	retList, err := utils.SqlQuery(odm.sqlCli, outboundDetailTable, &OutboundDetail{}, condition)
	if err != nil {
		logger.Warn(outboundDetailLogTag, "GetOutboundDetailByOrderList Failed|Condition:%v|outboundID:%#v|Err:%v",
			condition, outboundIDList, err)
		return nil, err
	}

	return retList.([]*OutboundDetail), nil
}
