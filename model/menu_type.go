package model

import (
	"database/sql"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
	"time"
)

const (
	menuTypeTable = "menu_type"

	menuTypeLogTag = "MenuTypeModel"
)

var (
	menuTypeUpdateTags = []string{"menu_config", "menu_type_name"}
)

type MenuType struct {
	ID           uint32    `json:"id"`
	MenuTypeName string    `json:"menu_type_name"`
	MenuConfig   string    `json:"menu_config"`
	CreateAt     time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}

type MenuTypeModel struct {
	sqlCli *sql.DB
}

func NewMenuTypeModelWithDB(sqlCli *sql.DB) *MenuTypeModel {
	return &MenuTypeModel{
		sqlCli: sqlCli,
	}
}

func (mtm *MenuTypeModel) Insert(dao *MenuType) error {
	id, err := utils.SqlInsert(mtm.sqlCli, menuTypeTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(menuTypeLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (mtm *MenuTypeModel) GetMenuTypes() ([]*MenuType, error) {
	retList, err := utils.SqlQuery(mtm.sqlCli, menuTypeTable, &MenuType{}, "")
	if err != nil {
		logger.Warn(menuTypeLogTag, "GetMenuTypes Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*MenuType), nil
}

func (mtm *MenuTypeModel) UpdateMenuType(dao *MenuType) error {
	err := utils.SqlUpdateWithUpdateTags(mtm.sqlCli, menuTypeTable, dao, "id", menuTypeUpdateTags...)
	if err != nil {
		logger.Warn(menuTypeLogTag, "UpdateMenuType Failed|Err:%v", err)
		return err
	}
	return nil
}
