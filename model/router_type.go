package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	routerTypeTable = "router_type"

	routerTypeLogTag = "RouterTypeModel"
)

var (
	routerTypeUpdateTags = []string{"router_type_name", "sort_id"}
)

type RouterType struct {
	ID             uint32    `json:"id"`
	RouterTypeName string    `json:"router_type_name"`
	SortID         uint32    `json:"sort_id"`
	CreateAt       time.Time `json:"created_at"`
	UpdateAt       time.Time `json:"updated_at"`
}

type RouterTypeModel struct {
	sqlCli *sql.DB
}

func NewRouterTypeModel(sqlCli *sql.DB) *RouterTypeModel {
	return &RouterTypeModel{
		sqlCli: sqlCli,
	}
}

func (rtm *RouterTypeModel) Insert(dao *RouterType) error {
	id, err := utils.SqlInsert(rtm.sqlCli, routerTypeTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(routerTypeLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (rtm *RouterTypeModel) GetRouterTypes() ([]*RouterType, error) {
	retList, err := utils.SqlQuery(rtm.sqlCli, routerTypeTable, &RouterType{}, " ORDER BY `sort_id` ASC ")
	if err != nil {
		logger.Warn(routerTypeLogTag, "GetRouterTypes Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*RouterType), nil
}

func (rtm *RouterTypeModel) UpdateRouterType(dao *RouterType) error {
	err := utils.SqlUpdateWithUpdateTags(rtm.sqlCli, routerTypeTable, dao, "id", routerTypeUpdateTags...)
	if err != nil {
		logger.Warn(routerTypeLogTag, "UpdateStorehouseType Failed|Err:%v", err)
		return err
	}
	return nil
}
