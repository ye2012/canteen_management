package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	routerDetailTable = "router_detail"

	routerDetailLogTag = "RouterDetail"
)

var (
	routerDetailUpdateTags = []string{"router_sort_id", "router_type", "router_name", "router_path", "role"}
)

type RouterDetail struct {
	ID           uint32    `json:"id"`
	RouterType   uint32    `json:"router_type"`
	RouterName   string    `json:"router_name"`
	RouterPath   string    `json:"router_path"`
	RouterSortID uint32    `json:"router_sort_id"`
	Role         uint32    `json:"role"`
	CreateAt     time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}

type RouterDetailModel struct {
	sqlCli *sql.DB
}

func NewRouterDetailModel(sqlCli *sql.DB) *RouterDetailModel {
	return &RouterDetailModel{sqlCli: sqlCli}
}

func (rdm *RouterDetailModel) Insert(dao *RouterDetail) error {
	id, err := utils.SqlInsert(rdm.sqlCli, routerDetailTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(routerDetailLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (rdm *RouterDetailModel) BatchInsert(tx *sql.Tx, routerDetails []*RouterDetail) (err error) {
	if tx != nil {
		err = utils.SqlInsertBatch(tx, routerDetailTable, routerDetails, "id")
	} else {
		err = utils.SqlInsertBatch(rdm.sqlCli, routerDetailTable, routerDetails, "id")
	}
	if err != nil {
		logger.Warn(routerDetailLogTag, "BatchInsert Failed|RouterDetail:%+v|Err:%v", routerDetails, err)
		return err
	}
	return nil
}

func (rdm *RouterDetailModel) GetRouterDetailByCondition(condition string, params ...interface{}) ([]*RouterDetail, error) {
	retList, err := utils.SqlQuery(rdm.sqlCli, routerDetailTable, &RouterDetail{}, condition, params...)
	if err != nil {
		logger.Warn(routerDetailLogTag, "GetRouterDetail Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*RouterDetail), nil
}

func (rdm *RouterDetailModel) GetRouterDetail(routerType uint32) ([]*RouterDetail, error) {
	params := make([]interface{}, 0)
	condition := ""
	if routerType > 0 {
		condition = " WHERE `router_type` = ? "
		params = append(params, routerType)
	}
	condition += " ORDER BY `router_sort_id` ASC "

	retList, err := utils.SqlQuery(rdm.sqlCli, routerDetailTable, &RouterDetail{}, condition, params...)
	if err != nil {
		logger.Warn(routerDetailLogTag, "GetRouterDetail Failed|Condition:%v|RouterType:%#v|Err:%v",
			condition, routerType, err)
		return nil, err
	}

	return retList.([]*RouterDetail), nil
}

func (rdm *RouterDetailModel) UpdateDetail(dao *RouterDetail) error {
	err := utils.SqlUpdateWithUpdateTags(rdm.sqlCli, routerDetailTable, dao, "id", routerDetailUpdateTags...)
	if err != nil {
		logger.Warn(routerDetailLogTag, "UpdateDetail Failed|CartID:%#v|Err:%v", dao.ID, err)
		return err
	}

	return nil
}

func (rdm *RouterDetailModel) DeleteWithTx(tx *sql.Tx, routerTypes []uint32) (err error) {
	if len(routerTypes) == 0 {
		return fmt.Errorf("cart id must be set|id:%v", routerTypes)
	}
	idStr := ""
	for _, cartID := range routerTypes {
		idStr += fmt.Sprintf(",%v", cartID)
	}
	condition := fmt.Sprintf(" WHERE `router_type` in (%v) ", idStr[1:])
	sqlStr := fmt.Sprintf("DELETE FROM `%v` %v ", routerDetailTable, condition)
	if tx != nil {
		_, err = tx.Exec(sqlStr)
	} else {
		_, err = rdm.sqlCli.Exec(sqlStr)
	}

	if err != nil {
		logger.Warn(routerDetailLogTag, "Delete Failed|Err:%v", err)
		return err
	}
	return nil
}

func (rdm *RouterDetailModel) Delete(routerTypes []uint32) error {
	return rdm.DeleteWithTx(nil, routerTypes)
}

func (rdm *RouterDetailModel) DeleteByID(routerID uint32) error {
	sqlStr := fmt.Sprintf(" DELETE FROM %v WHERE `id` = ? ", routerDetailTable)
	_, err := rdm.sqlCli.Exec(sqlStr, routerID)
	if err != nil {
		logger.Warn(routerDetailLogTag, "DeleteByID Failed|Err:%v", err)
		return err
	}

	return nil
}
