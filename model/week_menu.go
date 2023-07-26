package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	weekMenuTable = "week_menu"

	weekMenuLogTag = "WeekMenu"
)

var weekMenuUpdateTags = []string{"menu_content"}

type WeekMenu struct {
	ID            uint32    `json:"id"`
	MenuTypeID    uint32    `json:"menu_type_id"`
	MenuContent   string    `json:"menu_content"`
	MenuStartDate time.Time `json:"menu_start_date"`
	CreateAt      time.Time `json:"created_at"`
	UpdateAt      time.Time `json:"updated_at"`
}

type WeekMenuModel struct {
	sqlCli *sql.DB
}

func NewWeekMenuModelWithDB(sqlCli *sql.DB) *WeekMenuModel {
	return &WeekMenuModel{
		sqlCli: sqlCli,
	}
}

func (wmm *WeekMenuModel) Insert(dao *WeekMenu) error {
	id, err := utils.SqlInsert(wmm.sqlCli, weekMenuTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(weekMenuLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (wmm *WeekMenuModel) BatchInsert(menuList []*WeekMenu) error {
	err := utils.SqlInsertBatch(wmm.sqlCli, weekMenuTable, menuList, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(weekMenuLogTag, "Insert Failed|MenuList:%+v|Err:%v", menuList, err)
		return err
	}
	return nil
}

func (wmm *WeekMenuModel) GetWeekMenus(menuID, menuType uint32, startTime, endTime int64) ([]*WeekMenu, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if menuID > 0 {
		condition += " AND `id` = ? "
		params = append(params, menuID)
	}
	if menuType > 0 {
		condition += " AND `menu_type_id` = ? "
		params = append(params, menuType)
	}
	if startTime > 0 {
		start := time.Unix(startTime, 0)
		condition += " AND `menu_start_date` >= ? "
		params = append(params, start)
	}
	if endTime > startTime {
		end := time.Unix(endTime, 0)
		condition += "AND `menu_start_date` <= ? "
		params = append(params, end)
	}

	retList, err := utils.SqlQuery(wmm.sqlCli, weekMenuTable, &WeekMenu{}, condition, params...)
	if err != nil {
		logger.Warn(weekMenuLogTag, "GetMenus Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*WeekMenu), nil
}

func (wmm *WeekMenuModel) UpdateWeekMenu(dao *WeekMenu) error {
	err := utils.SqlUpdateWithUpdateTags(wmm.sqlCli, weekMenuTable, dao, "id", weekMenuUpdateTags...)
	if err != nil {
		logger.Warn(weekMenuLogTag, "UpdateMenu Failed|Err:%v", err)
		return err
	}
	return nil
}
