package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func (wm *WeekMenu) FromWeekMenuConfig(menuConf []map[uint8][]uint32) error {
	contentStr, err := json.Marshal(menuConf)
	if err != nil {
		logger.Warn(menuLogTag, "FromWeekMenuConfig Failed|Err:%v", err)
		return err
	}
	wm.MenuContent = string(contentStr)
	return nil
}

func (wm *WeekMenu) ToWeekMenuConfig() []map[uint8][]uint32 {
	configMap := make([]map[uint8][]uint32, 0)
	err := json.Unmarshal([]byte(wm.MenuContent), &configMap)
	if err != nil {
		logger.Warn(weekMenuLogTag, "ToWeekMenuConfig Failed|Err:%v", err)
		return nil
	}
	return configMap
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
		logger.Warn(weekMenuLogTag, "GetWeekMenus Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*WeekMenu), nil
}

func (wmm *WeekMenuModel) UpdateWeekMenu(dao *WeekMenu) error {
	err := utils.SqlUpdateWithUpdateTags(wmm.sqlCli, weekMenuTable, dao, "id", weekMenuUpdateTags...)
	if err != nil {
		logger.Warn(weekMenuLogTag, "UpdateWeekMenu Failed|Err:%v", err)
		return err
	}
	return nil
}

func (wmm *WeekMenuModel) GetWeekMenuByDate(menuDate int64, menuType uint32) (*WeekMenu, error) {
	condition := " WHERE `menu_start_date` = ? AND `menu_type_id` = ? "
	retList, err := utils.SqlQuery(wmm.sqlCli, weekMenuTable, &WeekMenu{}, condition, time.Unix(menuDate, 0), menuType)
	if err != nil {
		logger.Warn(weekMenuLogTag, "GetWeekMenuByDate Failed|Err:%v", err)
		return nil, err
	}

	menuList := retList.([]*WeekMenu)
	if len(menuList) == 0 {
		logger.Warn(weekMenuLogTag, "WeekMenu Not Found|Date:%v|Type:%v", time.Unix(menuDate, 0), menuType)
		return nil, fmt.Errorf("week menu not found")
	}

	return menuList[0], nil
}
