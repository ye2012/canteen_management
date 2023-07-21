package service

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	menuServiceLogTag = "MenuService"
)

type MenuService struct {
	menuModel     *model.MenuModel
	menuTypeModel *model.MenuTypeModel

	menuTypeMap map[uint32]*model.MenuType
}

func NewMenuService(sqlCli *sql.DB) *MenuService {
	menuModel := model.NewMenuModelWithDB(sqlCli)
	menuTypeModel := model.NewMenuTypeModelWithDB(sqlCli)
	return &MenuService{
		menuModel:     menuModel,
		menuTypeModel: menuTypeModel,
		menuTypeMap:   make(map[uint32]*model.MenuType),
	}
}

func (ms *MenuService) Init() error {
	typeList, err := ms.GetMenuTypeList()
	if err != nil {
		logger.Warn(menuServiceLogTag, "Init Failed|Err:%v", err)
		return err
	}

	for _, typeInfo := range typeList {
		ms.menuTypeMap[typeInfo.ID] = typeInfo
	}
	return nil
}

func (ms *MenuService) GetMenuTypeList() ([]*model.MenuType, error) {
	typeList, err := ms.menuTypeModel.GetMenuTypes()
	if err != nil {
		logger.Warn(menuServiceLogTag, "GetMenuTypeList Failed|Err:%v", err)
		return nil, err
	}
	return typeList, nil
}

func (ms *MenuService) GetMenuTypeConfig(typeID uint32) *model.MenuType {
	_, ok := ms.menuTypeMap[typeID]
	if ok {
		return ms.menuTypeMap[typeID]
	}
	return nil
}

func (ms *MenuService) AddMenuType(dao *model.MenuType) error {
	err := ms.menuTypeModel.Insert(dao)
	if err != nil {
		logger.Warn(menuServiceLogTag, "Insert MenuType Failed|Err:%v", err)
		return err
	}
	ms.updateMenuTypeCache(dao)
	return nil
}

func (ms *MenuService) UpdateMenuType(dao *model.MenuType) error {
	err := ms.menuTypeModel.UpdateMenuType(dao)
	if err != nil {
		logger.Warn(menuServiceLogTag, "Update MenuType Failed|Err:%v", err)
		return err
	}
	ms.updateMenuTypeCache(dao)
	return nil
}

func (ms *MenuService) updateMenuTypeCache(dao *model.MenuType) {
	ms.menuTypeMap[dao.ID] = dao
}

func (ms *MenuService) GetMenuList(menuType uint32, startTime, endTime int64) ([]*model.Menu, error) {
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)
	menuList, err := ms.menuModel.GetMenus(menuType, start, end)
	if err != nil {
		logger.Warn(menuServiceLogTag, "GetMenuList Failed|Err:%v", err)
		return nil, err
	}

	return menuList, nil
}

func (ms *MenuService) AddMenu(menu *model.Menu) error {
	err := ms.menuModel.Insert(menu)
	if err != nil {
		logger.Warn(menuServiceLogTag, "Insert Menu Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ms *MenuService) BatchAddMenu(menuList []*model.Menu) error {
	err := ms.menuModel.BatchInsert(menuList)
	if err != nil {
		logger.Warn(menuServiceLogTag, "BatchInsert Menu Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ms *MenuService) UpdateMenu(menu *model.Menu) error {
	err := ms.menuModel.UpdateMenu(menu)
	if err != nil {
		logger.Warn(menuServiceLogTag, "Update Menu Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ms *MenuService) GenerateMenu(menuType uint32) {

}
