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
	orderTable = "order"

	orderLogTag = "OrderModel"
)

type OrderDao struct {
	ID uint32 `json:"id"`
	//OrderID      string          `json:"order_id"`
	PayOrderID     uint32    `json:"pay_order_id"`
	MealType       uint8     `json:"meal_type"`
	OrderDate      time.Time `json:"order_date"`
	Uid            uint32    `json:"uid"`
	PhoneNumber    string    `json:"phone_number"`
	BuildingID     uint32    `json:"building_id"`
	Floor          uint32    `json:"floor"`
	Room           string    `json:"room"`
	TotalAmount    float64   `json:"total_amount"`
	PayMethod      uint8     `json:"pay_method"`
	PayAmount      float64   `json:"pay_amount"`
	DiscountAmount float64   `json:"discount_amount"`
	DeliverUid     uint32    `json:"deliver_user_id"`
	DeliverTime    time.Time `json:"deliver_time"`
	Status         uint8     `json:"status"`
	CreateAt       time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"updated_at"`
}

type OrderModel struct {
	sqlCli *sql.DB
}

func NewOrderModel(sqlCli *sql.DB) *OrderModel {
	return &OrderModel{
		sqlCli: sqlCli,
	}
}

func (om *OrderModel) InsertWithTx(tx *sql.Tx, dao *OrderDao) (err error) {
	id := int64(0)
	if tx != nil {
		id, err = utils.SqlInsert(tx, orderTable, dao, "id", "created_at", "updated_at", "deliver_time")
	} else {
		id, err = utils.SqlInsert(om.sqlCli, orderTable, dao, "id", "created_at", "updated_at", "deliver_time")
	}

	if err != nil {
		logger.Warn(orderLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (om *OrderModel) UpdateOrderInfoByID(tx *sql.Tx, order *OrderDao, updateTags ...string) (err error) {
	return om.UpdateOrderInfo(tx, order, "id", updateTags...)
}

func (om *OrderModel) UpdateOrderInfo(tx *sql.Tx, order *OrderDao, conditionTag string, updateTags ...string) (err error) {
	if tx != nil {
		err = utils.SqlUpdateWithUpdateTags(tx, orderTable, order, conditionTag, updateTags...)
	} else {
		err = utils.SqlUpdateWithUpdateTags(om.sqlCli, orderTable, order, conditionTag, updateTags...)
	}
	if err != nil {
		logger.Warn(orderLogTag, "UpdateOrderInfoByID Failed|Err:%v", err)
		return err
	}
	return nil
}

func (om *OrderModel) GetOrderListByPayOrder(payOrderList []uint32) ([]*OrderDao, error) {
	if len(payOrderList) == 0 {
		return make([]*OrderDao, 0), nil
	}
	payOrderStr := ""
	for _, payOrderID := range payOrderList {
		payOrderStr += fmt.Sprintf(",%v", payOrderID)
	}
	condition := fmt.Sprintf(" WHERE `pay_order_id` in (%v) ", payOrderStr[1:])
	retList, err := utils.SqlQuery(om.sqlCli, orderTable, &OrderDao{}, condition)
	if err != nil {
		logger.Warn(orderLogTag, "GetOrderListByPayOrder Failed|PayOrderID:%v|Err:%v", payOrderList, err)
		return nil, err
	}

	return retList.([]*OrderDao), nil
}

func (om *OrderModel) GenerateCondition(idList []uint32, uid uint32, status int8, buildingID, floor uint32,
	room string, startTime, endTime int64, mealType uint8, payMethod int8) (string, []interface{}) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if len(idList) > 0 {
		idStr := ""
		for _, id := range idList {
			idStr += fmt.Sprintf(",%v", id)
		}
		condition += fmt.Sprintf(" AND `id` in (%v) ", idStr)
	}
	if uid > 0 {
		condition += " AND `uid` = ? "
		params = append(params, uid)
	}
	if status != -1 {
		condition += " AND `status` = ? "
		params = append(params, status)
	}
	if buildingID > 0 {
		condition += " AND `building_id` = ? "
		params = append(params, buildingID)
	}
	if floor > 0 {
		condition += " AND `floor` = ? "
		params = append(params, floor)
	}
	if startTime > 0 {
		condition += " AND `order_date` >= ? "
		params = append(params, time.Unix(startTime, 0))
	}
	if endTime > startTime {
		condition += " AND `order_date` <= ? "
		params = append(params, time.Unix(endTime, 0))
	}
	if mealType > enum.MealUnknown {
		condition += " AND `meal_type` = ? "
		params = append(params, mealType)
	}
	if payMethod > -1 {
		condition += " AND `pay_method` = ? "
		params = append(params, payMethod)
	}
	if room != "" {
		condition += " AND `room` = ? "
		params = append(params, room)
	}
	return condition, params
}

func (om *OrderModel) GetFloors(buildingID uint32, status, payMethod int8, startTime, endTime int64, mealType uint8) ([]int32, error) {
	condition, params := om.GenerateCondition(make([]uint32, 0), 0, status, buildingID, 0, "",
		startTime, endTime, mealType, payMethod)
	sqlStr := fmt.Sprintf("SELECT DISTINCT(`floor`) FROM `%v` %v", orderTable, condition)
	rows, err := om.sqlCli.Query(sqlStr, params...)
	if err != nil {
		logger.Warn(orderLogTag, "GetFloors Query Failed|BuildingID:%v|Status:%v|MealType:%v|Err:%v",
			buildingID, status, mealType, err)
		return nil, err
	}
	defer rows.Close()

	floors := make([]int32, 0)
	for rows.Next() {
		floor := int32(0)
		err = rows.Scan(&floor)
		if err != nil {
			logger.Warn(orderLogTag, "GetFloors Scan Failed|Err:%v", err)
			continue
		}
		floors = append(floors, floor)
	}

	return floors, nil
}

func (om *OrderModel) GetAllOrder(mealType uint8, startTime, endTime int64, status, payMethod int8) ([]*OrderDao, error) {
	condition, params := om.GenerateCondition(make([]uint32, 0), 0, status, 0, 0, "",
		startTime, endTime, mealType, payMethod)
	retList, err := utils.SqlQuery(om.sqlCli, orderTable, &OrderDao{}, condition, params...)
	if err != nil {
		logger.Warn(orderLogTag, "GetAllOrder Failed|MealType:%v|Start:%v|End:%v|Status:%v|Err:%v",
			mealType, startTime, endTime, status, err)
		return nil, err
	}

	return retList.([]*OrderDao), nil
}

func (om *OrderModel) GetOrderList(idList []uint32, uid uint32, mealType uint8, buildingID, floor uint32, room string,
	status, payMethod int8, startTime, endTime int64, page, pageSize int32) ([]*OrderDao, error) {
	condition, params := om.GenerateCondition(idList, uid, status, buildingID, floor, room, startTime, endTime,
		mealType, payMethod)
	condition += " ORDER BY `id` ASC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(om.sqlCli, orderTable, &OrderDao{}, condition, params...)
	if err != nil {
		logger.Warn(orderLogTag, "GetOrderList Failed|ID:%v|Uid:%v|Status:%v|Err:%v", idList, uid, status, err)
		return nil, err
	}

	return retList.([]*OrderDao), nil
}

func (om *OrderModel) GetOrderListCount(idList []uint32, uid uint32, mealType uint8, status, payMethod int8, buildingID,
	floor uint32, room string, startTime, endTime int64) (int32, error) {
	condition, params := om.GenerateCondition(idList, uid, status, buildingID, floor, room,
		startTime, endTime, mealType, payMethod)

	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM `%v` %v", orderTable, condition)
	row := om.sqlCli.QueryRow(sqlStr, params...)
	var count int32 = 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (om *OrderModel) GetOrder(id uint32) (*OrderDao, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	retInfo := &OrderDao{}
	err := utils.SqlQueryRow(om.sqlCli, orderTable, retInfo, condition, params...)
	if err != nil {
		logger.Warn(orderLogTag, "GetOrder Failed|ID:%v|Err:%v", id, err)
		return nil, err
	}

	return retInfo, nil
}
