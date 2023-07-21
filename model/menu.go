package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	menuTable = "menu"

	menuLogTag = "Menu"
)

var menuUpdateTags = []string{"menu_content", ""}

type Menu struct {
	ID          uint32    `json:"id"`
	MenuTypeID  uint32    `json:"menu_type_id"`
	MenuContent string    `json:"menu_content"`
	MenuDate    time.Time `json:"menu_date"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

type MenuModel struct {
	sqlCli *sql.DB
}

func NewMenuModelWithDB(sqlCli *sql.DB) *MenuModel {
	return &MenuModel{
		sqlCli: sqlCli,
	}
}

func (mm *MenuModel) Insert(dao *Menu) error {
	id, err := utils.SqlInsert(mm.sqlCli, menuTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(menuLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (mm *MenuModel) BatchInsert(menuList []*Menu) error {
	err := utils.SqlInsertBatch(mm.sqlCli, menuTable, menuList, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(menuLogTag, "Insert Failed|MenuList:%+v|Err:%v", menuList, err)
		return err
	}
	return nil
}

func (mm *MenuModel) GetMenus(menuType uint32, startTime, endTime time.Time) ([]*Menu, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if menuType > 0 {
		condition += " AND `menu_type_id` = ? "
		params = append(params, menuType)
	}
	condition += " AND `menu_date` >= ? AND `menu_date` <= ? "
	params = append(params, startTime, endTime)
	retList, err := utils.SqlQuery(mm.sqlCli, menuTable, &Menu{}, condition, params...)
	if err != nil {
		logger.Warn(menuLogTag, "GetMenus Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*Menu), nil
}

func (mm *MenuModel) UpdateMenu(dao *Menu) error {
	err := utils.SqlUpdateWithUpdateTags(mm.sqlCli, menuTable, dao, "id", menuUpdateTags...)
	if err != nil {
		logger.Warn(menuLogTag, "UpdateMenu Failed|Err:%v", err)
		return err
	}
	return nil
}
