package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	payOrderTable = "pay_order"

	payOrderLogTag = "PayOrderModel"
)

var (
	payOrderUpdateTags = []string{"status"}
)

type PayOrderDao struct {
	ID           uint32    `json:"id"`
	PrepareID    string    `json:"prepare_id"`
	Uid          uint32    `json:"uid"`
	UnionID      string    `json:"union_id"`
	Address      string    `json:"address"`
	TotalAmount  float64   `json:"total_amount"`
	PayAmount    float64   `json:"pay_amount"`
	DiscountType uint8     `json:"discount_type"`
	Status       uint8     `json:"status"`
	CreateAt     time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}

type PayOrderModel struct {
	sqlCli *sql.DB
}

func NewPayOrderModel(sqlCli *sql.DB) *PayOrderModel {
	return &PayOrderModel{
		sqlCli: sqlCli,
	}
}

func (pom *PayOrderModel) InsertWithTx(tx *sql.Tx, dao *PayOrderDao) (err error) {
	id := int64(0)
	if tx != nil {
		id, err = utils.SqlInsert(tx, payOrderTable, dao, "id", "created_at", "updated_at")
	} else {
		id, err = utils.SqlInsert(pom.sqlCli, payOrderTable, dao, "id", "created_at", "updated_at")
	}

	if err != nil {
		logger.Warn(payOrderLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (pom *PayOrderModel) UpdatePayOrderInfoByID(tx *sql.Tx, order *PayOrderDao, updateTags ...string) (err error) {
	if tx != nil {
		err = utils.SqlUpdateWithUpdateTags(tx, payOrderTable, order, "id", updateTags...)
	} else {
		err = utils.SqlUpdateWithUpdateTags(pom.sqlCli, payOrderTable, order, "id", updateTags...)
	}
	if err != nil {
		logger.Warn(payOrderLogTag, "UpdateOrderInfoByID Failed|Err:%v", err)
		return err
	}
	return nil
}

func (pom *PayOrderModel) GenerateCondition(idList []uint32, uid uint32, page, pageSize int32, status int8) (string, []interface{}) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if len(idList) > 0 {
		idStr := ""
		for _, id := range idList {
			idStr += fmt.Sprintf("%v", id)
		}
		condition += fmt.Sprintf(" AND `id` in (%v) ", idList[1:])
	}
	if uid > 0 {
		condition += " AND `uid` = ? "
		params = append(params, uid)
	}
	if status != -1 {
		condition += " AND `status` = ? "
		params = append(params, status)
	}
	condition += " ORDER BY `id` ASC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	return condition, params
}

func (pom *PayOrderModel) GetPayOrderList(idList []uint32, uid uint32, page, pageSize int32, status int8) ([]*PayOrderDao, error) {
	condition, params := pom.GenerateCondition(idList, uid, page, pageSize, status)
	retList, err := utils.SqlQuery(pom.sqlCli, payOrderTable, &PayOrderDao{}, condition, params...)
	if err != nil {
		logger.Warn(payOrderLogTag, "GetPayOrderList Failed|ID:%v|Uid:%v|Status:%v|Err:%v", idList, uid, status, err)
		return nil, err
	}

	return retList.([]*PayOrderDao), nil
}

func (pom *PayOrderModel) GetPayOrderListCount(idList []uint32, uid uint32, page, pageSize int32, status int8) (int32, error) {
	condition, params := pom.GenerateCondition(idList, uid, page, pageSize, status)
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM %v %v", payOrderTable, condition)
	row := pom.sqlCli.QueryRow(sqlStr, params...)
	var count int32 = 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pom *PayOrderModel) GetPayOrder(id uint32) (*PayOrderDao, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	retInfo := &PayOrderDao{}
	err := utils.SqlQueryRow(pom.sqlCli, payOrderTable, retInfo, condition, params...)
	if err != nil {
		logger.Warn(payOrderLogTag, "GetOrder Failed|ID:%v|Err:%v", id, err)
		return nil, err
	}

	return retInfo, nil
}
