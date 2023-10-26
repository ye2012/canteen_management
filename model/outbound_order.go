package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	outboundOrderTable = "outbound_order"

	outboundOrderLogTag = "OutboundModel"
)

type OutboundOrder struct {
	ID          uint32    `json:"id"`
	Creator     uint32    `json:"creator"`
	TotalAmount float64   `json:"total_amount"`
	Status      int8      `json:"status"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

type OutboundOrderModel struct {
	sqlCli *sql.DB
}

func NewOutboundOrderModelWithDB(sqlCli *sql.DB) *OutboundOrderModel {
	return &OutboundOrderModel{
		sqlCli: sqlCli,
	}
}

func (oom *OutboundOrderModel) Insert(dao *OutboundOrder) error {
	return oom.InsertWithTx(nil, dao)
}

func (oom *OutboundOrderModel) InsertWithTx(tx *sql.Tx, dao *OutboundOrder) error {
	id, err := int64(0), error(nil)
	if tx != nil {
		id, err = utils.SqlInsert(tx, outboundOrderTable, dao, "id", "created_at", "updated_at")
	} else {
		id, err = utils.SqlInsert(oom.sqlCli, outboundOrderTable, dao, "id", "created_at", "updated_at")
	}
	if err != nil {
		logger.Warn(outboundOrderLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (oom *OutboundOrderModel) GenerateCondition(id, creator uint32, startTime, endTime int64, status int8) (string, []interface{}) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	if creator > 0 {
		condition += " AND `creator` = ? "
		params = append(params, creator)
	}
	if startTime > 0 {
		condition += " AND `created_at` >= ? "
		params = append(params, time.Unix(startTime, 0))
	}
	if endTime > startTime {
		condition += " AND `created_at` <= ? "
		params = append(params, time.Unix(endTime, 0))
	}
	if status != -1 {
		condition += " AND `status` = ? "
		params = append(params, status)
	}
	return condition, params
}

func (oom *OutboundOrderModel) GetOutboundOrder(id uint32) (*OutboundOrder, error) {
	condition := " WHERE 1=1 "
	params := make([]interface{}, 0)
	if id > 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	retInfo := &OutboundOrder{}
	err := utils.SqlQueryRow(oom.sqlCli, outboundOrderTable, retInfo, condition, params...)
	if err != nil {
		logger.Warn(outboundOrderLogTag, "GetOrder Failed|ID:%v|Err:%v", id, err)
		return nil, err
	}

	return retInfo, nil
}

func (oom *OutboundOrderModel) GetOutboundOrderWithLock(tx *sql.Tx, id uint32) (*OutboundOrder, error) {
	if tx == nil {
		return nil, fmt.Errorf("tx is nil")
	}
	condition := " WHERE `id` = ? "
	retInfo := &OutboundOrder{}
	err := utils.SqlQueryRowWithLock(oom.sqlCli, outboundOrderTable, retInfo, condition, id)
	if err != nil {
		logger.Warn(outboundOrderLogTag, "GetOrderWithLock Failed|ID:%v|Err:%v", id, err)
		return nil, err
	}

	return retInfo, nil
}

func (oom *OutboundOrderModel) GetOutboundOrderList(id uint32, creator uint32, startTime, endTime int64, status int8,
	page, pageSize int32) ([]*OutboundOrder, error) {
	condition, params := oom.GenerateCondition(id, creator, startTime, endTime, status)
	condition += " ORDER BY `id` DESC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(oom.sqlCli, outboundOrderTable, &OutboundOrder{}, condition, params...)
	if err != nil {
		logger.Warn(outboundOrderLogTag, "GetOutboundOrders Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*OutboundOrder), nil
}

func (oom *OutboundOrderModel) GetOutboundOrderCount(id uint32, creator uint32, startTime,
	endTime int64, status int8) (int32, error) {
	condition, params := oom.GenerateCondition(id, creator, startTime, endTime, status)
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM `%v` %v", outboundOrderTable, condition)
	row := oom.sqlCli.QueryRow(sqlStr, params...)
	var count int32 = 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (oom *OutboundOrderModel) UpdateOutboundWithTx(tx *sql.Tx, dao *OutboundOrder, updateTags ...string) (err error) {
	if tx != nil {
		err = utils.SqlUpdateWithUpdateTags(tx, outboundOrderTable, dao, "id", updateTags...)
	} else {
		err = utils.SqlUpdateWithUpdateTags(oom.sqlCli, outboundOrderTable, dao, "id", updateTags...)
	}
	if err != nil {
		logger.Warn(outboundOrderLogTag, "UpdateOutboundWithTx Failed|Err:%v", err)
		return err
	}
	return nil
}
