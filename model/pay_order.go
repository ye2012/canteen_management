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
	ID             uint32    `json:"id"`
	PrepareID      string    `json:"prepare_id"`
	MealTime       time.Time `json:"meal_time"`
	Uid            uint32    `json:"uid"`
	OpenID         string    `json:"open_id"`
	CartID         uint32    `json:"cart_id"`
	BuildingID     uint32    `json:"building_id"`
	Floor          uint32    `json:"floor"`
	Room           string    `json:"room"`
	TotalAmount    float64   `json:"total_amount"`
	PayMethod      uint8     `json:"pay_method"`
	PayAmount      float64   `json:"pay_amount"`
	DiscountAmount float64   `json:"discount_amount"`
	Status         uint8     `json:"status"`
	CreateAt       time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"updated_at"`
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

func (pom *PayOrderModel) GenerateCondition(idList []uint32, uid, cartID uint32, statusList []int8, timeStart, timeEnd int64) (string, []interface{}) {
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
	if cartID > 0 {
		condition += " AND `cart_id` = ? "
		params = append(params, cartID)
	}
	if len(statusList) > 0 {
		statusStr := ""
		for _, status := range statusList {
			statusStr += fmt.Sprintf(",%v", status)
		}
		condition += fmt.Sprintf(" AND `status` in (%v) ", statusStr[1:])
		params = append(params)
	}
	if timeStart > 0 {
		condition += " AND `meal_time` >= ? "
		params = append(params, time.Unix(timeStart, 0))
	}
	if timeEnd > 0 {
		condition += " AND `meal_time` <= ? "
		params = append(params, time.Unix(timeEnd, 0))
	}
	return condition, params
}

func (pom *PayOrderModel) GetPayOrderList(idList []uint32, uid uint32, page, pageSize int32, status int8) ([]*PayOrderDao, error) {
	statusList := make([]int8, 0)
	if status != -1 {
		statusList = append(statusList, status)
	}
	condition, params := pom.GenerateCondition(idList, uid, 0, statusList, 0, 0)
	condition += " ORDER BY `id` ASC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(pom.sqlCli, payOrderTable, &PayOrderDao{}, condition, params...)
	if err != nil {
		logger.Warn(payOrderLogTag, "GetPayOrderList Failed|ID:%v|Uid:%v|Status:%v|Err:%v", idList, uid, status, err)
		return nil, err
	}

	return retList.([]*PayOrderDao), nil
}

func (pom *PayOrderModel) GetAllPayOrderList(idList []uint32, uid uint32, statusList []int8,
	timeStart, timeEnd int64) ([]*PayOrderDao, error) {
	condition, params := pom.GenerateCondition(idList, uid, 0, statusList, timeStart, timeEnd)
	retList, err := utils.SqlQueryWithLock(pom.sqlCli, payOrderTable, &PayOrderDao{}, condition, params...)
	if err != nil {
		logger.Warn(payOrderLogTag, "GetPayOrderList Failed|ID:%v|Uid:%v|Status:%v|Err:%v", idList, uid, statusList, err)
		return nil, err
	}

	return retList.([]*PayOrderDao), nil
}

func (pom *PayOrderModel) GetPayOrderListWithLock(tx *sql.Tx, idList []uint32, uid uint32, statusList []int8,
	timeStart, timeEnd int64) ([]*PayOrderDao, error) {
	if tx == nil {
		return nil, fmt.Errorf("tx can not be null")
	}
	condition, params := pom.GenerateCondition(idList, uid, 0, statusList, timeStart, timeEnd)
	retList, err := utils.SqlQueryWithLock(tx, payOrderTable, &PayOrderDao{}, condition, params...)
	if err != nil {
		logger.Warn(payOrderLogTag, "GetPayOrderListWithLock Failed|ID:%v|Uid:%v|Status:%v|Err:%v",
			idList, uid, statusList, err)
		return nil, err
	}

	return retList.([]*PayOrderDao), nil
}

func (pom *PayOrderModel) GetPayOrderListCount(idList []uint32, uid uint32, status int8) (int32, error) {
	statusList := make([]int8, 0)
	if status != -1 {
		statusList = append(statusList, status)
	}
	condition, params := pom.GenerateCondition(idList, uid, 0, statusList, 0, 0)
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM %v %v", payOrderTable, condition)
	row := pom.sqlCli.QueryRow(sqlStr, params...)
	var count int32 = 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pom *PayOrderModel) GetPayOrderByCondition(tx *sql.Tx, condition string, params ...interface{}) (*PayOrderDao, error) {
	retInfo, err := &PayOrderDao{}, error(nil)
	if tx != nil {
		err = utils.SqlQueryRowWithLock(tx, payOrderTable, retInfo, condition, params...)
	} else {
		err = utils.SqlQueryRow(pom.sqlCli, payOrderTable, retInfo, condition, params...)
	}
	if err != nil {
		logger.Warn(payOrderLogTag, "GetPayOrderByCondition Failed|Condition:%v|params:%v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retInfo, nil
}

func (pom *PayOrderModel) GetPayOrder(id uint32) (*PayOrderDao, error) {
	return pom.GetPayOrderByCondition(nil, " WHERE `id` = ? ", id)
}
