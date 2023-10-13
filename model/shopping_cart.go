package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	shoppingCartTable = "shopping_cart"

	shoppingCartLogTag = "ShoppingCart"
)

type ShoppingCart struct {
	ID       uint32    `json:"id"`
	CartType uint8     `json:"cart_type"`
	Uid      uint32    `json:"uid"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}
type ShoppingCartModel struct {
	sqlCli *sql.DB
}

func NewShoppingCartModel(sqlCli *sql.DB) *ShoppingCartModel {
	return &ShoppingCartModel{sqlCli: sqlCli}
}

func (scm *ShoppingCartModel) InsertWithTx(tx *sql.Tx, dao *ShoppingCart) (err error) {
	id := int64(0)
	if tx != nil {
		id, err = utils.SqlInsert(tx, shoppingCartTable, dao, "id", "created_at", "updated_at")
	} else {
		id, err = utils.SqlInsert(scm.sqlCli, shoppingCartTable, dao, "id", "created_at", "updated_at")
	}

	if err != nil {
		logger.Warn(shoppingCartLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (scm *ShoppingCartModel) GenerateCondition(cartType uint8, uid, cartID uint32) (string, []interface{}) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if cartType != 0 {
		condition += " AND `cart_type` = ? "
		params = append(params, cartType)
	}
	if uid != 0 {
		condition += " AND `uid` = ? "
		params = append(params, uid)
	}
	if cartID != 0 {
		condition += " AND `id` = ? "
		params = append(params, cartID)
	}
	return condition, params
}

func (scm *ShoppingCartModel) GetCart(cartType uint8, uid uint32) ([]*ShoppingCart, error) {
	condition, params := scm.GenerateCondition(cartType, uid, 0)
	condition += " ORDER BY `created_at` DESC "
	retList, err := utils.SqlQuery(scm.sqlCli, shoppingCartTable, &ShoppingCart{}, condition, params...)
	if err != nil {
		logger.Warn(shoppingCartLogTag, "GetCart Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*ShoppingCart), nil
}

func (scm *ShoppingCartModel) GetCartByID(cartID uint32) (*ShoppingCart, error) {
	condition, params := scm.GenerateCondition(0, 0, cartID)
	cart := &ShoppingCart{}
	err := utils.SqlQueryRow(scm.sqlCli, shoppingCartTable, cart, condition, params...)
	if err != nil {
		logger.Warn(shoppingCartLogTag, "GetCartByID Failed|Err:%v", err)
		return nil, err
	}

	return cart, nil
}

func (scm *ShoppingCartModel) GetCartByTxWithLock(tx *sql.Tx, cartType uint8, uid uint32) ([]*ShoppingCart, error) {
	condition, params := scm.GenerateCondition(cartType, uid, 0)
	condition += " ORDER BY `created_at` DESC "
	retList := make([]*ShoppingCart, 0)
	if tx != nil {
		ret, err := utils.SqlQueryWithLock(tx, shoppingCartTable, &ShoppingCart{}, condition, params...)
		if err != nil {
			logger.Warn(shoppingCartLogTag, "GetCardByTx Failed|Err:%v", err)
			return nil, err
		}
		retList = ret.([]*ShoppingCart)
	} else {
		ret, err := utils.SqlQueryWithLock(scm.sqlCli, shoppingCartTable, &ShoppingCart{}, condition, params...)
		if err != nil {
			logger.Warn(shoppingCartLogTag, "GetCardByTx Failed|Err:%v", err)
			return nil, err
		}
		retList = ret.([]*ShoppingCart)
	}

	return retList, nil
}

func (scm *ShoppingCartModel) DeleteWithTx(tx *sql.Tx, cartType uint8, uid, cartID uint32) (err error) {
	condition, params := scm.GenerateCondition(cartType, uid, cartID)
	sqlStr := fmt.Sprintf("DELETE FROM `%v` %v ", shoppingCartTable, condition)
	if tx != nil {
		_, err = tx.Exec(sqlStr, params...)
	} else {
		_, err = scm.sqlCli.Exec(sqlStr, params...)
	}

	if err != nil {
		logger.Warn(shoppingCartLogTag, "Delete Failed|Err:%v", err)
		return err
	}
	return nil
}

func (scm *ShoppingCartModel) Delete(cartType uint8, uid, cartID uint32) error {
	return scm.DeleteWithTx(nil, cartType, uid, cartID)
}

func (scm *ShoppingCartModel) DeleteByTx(tx *sql.Tx, cartType uint8, uid, cartID uint32) error {
	return scm.DeleteWithTx(tx, cartType, uid, cartID)
}
