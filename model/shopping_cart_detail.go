package model

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	cartDetailTable = "shopping_cart_detail"

	cartDetailLogTag = "CartDetail"
)

var (
	cartDetailUpdateTags = []string{"quantity"}
)

type CartDetail struct {
	ID       uint32  `json:"id"`
	CartID   uint32  `json:"cart_id"`
	ItemID   string  `json:"item_id"`
	Quantity float64 `json:"quantity"`
}

type CartDetailModel struct {
	sqlCli *sql.DB
}

func NewCartDetailModel(sqlCli *sql.DB) *CartDetailModel {
	return &CartDetailModel{sqlCli: sqlCli}
}

func (scd *CartDetailModel) BatchInsert(tx *sql.Tx, shoppingDetail []*CartDetail) (err error) {
	if tx != nil {
		err = utils.SqlInsertBatch(tx, cartDetailTable, shoppingDetail, "id")
	} else {
		err = utils.SqlInsertBatch(scd.sqlCli, cartDetailTable, shoppingDetail, "id")
	}
	if err != nil {
		logger.Warn(cartDetailLogTag, "BatchInsert Failed|CartDetail:%+v|Err:%v", shoppingDetail, err)
		return err
	}
	return nil
}

func (scd *CartDetailModel) GetCartDetailByCondition(condition string, params ...interface{}) ([]*CartDetail, error) {
	retList, err := utils.SqlQuery(scd.sqlCli, cartDetailTable, &CartDetail{}, condition, params...)
	if err != nil {
		logger.Warn(cartDetailLogTag, "GetCartDetail Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*CartDetail), nil
}

func (scd *CartDetailModel) GetCartDetail(cartID uint32) ([]*CartDetail, error) {
	condition := " WHERE `cart_id` = ? "
	retList, err := utils.SqlQuery(scd.sqlCli, cartDetailTable, &CartDetail{}, condition, cartID)
	if err != nil {
		logger.Warn(cartDetailLogTag, "GetCartDetail Failed|Condition:%v|CartID:%#v|Err:%v",
			condition, cartID, err)
		return nil, err
	}

	return retList.([]*CartDetail), nil
}

func (scd *CartDetailModel) UpdateDetail(tx *sql.Tx, dao *CartDetail) error {
	if tx == nil {
		return fmt.Errorf("tx is nil")
	}
	err := utils.SqlUpdateWithUpdateTags(tx, cartDetailTable, dao, "id", cartDetailUpdateTags...)
	if err != nil {
		logger.Warn(cartDetailLogTag, "UpdateDetail Failed|CartID:%#v|Err:%v", dao.ID, err)
		return err
	}

	return nil
}

func (scd *CartDetailModel) GetCartDetailWithLock(tx *sql.Tx, cartID uint32) ([]*CartDetail, error) {
	condition := " WHERE `cart_id` = ? "
	if tx == nil {
		return nil, fmt.Errorf("tx nil")
	}
	retList, err := utils.SqlQueryWithLock(tx, cartDetailTable, &CartDetail{}, condition, cartID)
	if err != nil {
		logger.Warn(cartDetailLogTag, "GetCartDetail Failed|Condition:%v|CartID:%#v|Err:%v",
			condition, cartID, err)
		return nil, err
	}

	return retList.([]*CartDetail), nil
}

func (scd *CartDetailModel) DeleteWithTx(tx *sql.Tx, cartIDs []uint32) (err error) {
	if len(cartIDs) == 0 {
		return fmt.Errorf("cart id must be set|id:%v", cartIDs)
	}
	idStr := ""
	for _, cartID := range cartIDs {
		idStr += fmt.Sprintf(",%v", cartID)
	}
	condition := fmt.Sprintf(" WHERE `cart_id` in (%v) ", idStr[1:])
	sqlStr := fmt.Sprintf("DELETE FROM `%v` %v ", cartDetailTable, condition)
	if tx != nil {
		_, err = tx.Exec(sqlStr)
	} else {
		_, err = scd.sqlCli.Exec(sqlStr)
	}

	if err != nil {
		logger.Warn(cartDetailLogTag, "Delete Failed|Err:%v", err)
		return err
	}
	return nil
}

func (scd *CartDetailModel) Delete(cartIDs []uint32) error {
	return scd.DeleteWithTx(nil, cartIDs)
}
