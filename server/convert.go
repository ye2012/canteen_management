package server

import (
	"encoding/json"
	"time"

	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	convertLogTag = "Convert"
)

func ConvertToDishTypeInfoList(daoList []*model.DishType) []*dto.DishTypeInfo {
	retList := make([]*dto.DishTypeInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.DishTypeInfo{DishTypeID: dao.ID, DishTypeName: dao.DishTypeName,
			MasterTypeID: dao.MasterType})
	}
	return retList
}

func ConvertFromDishTypeInfo(info *dto.DishTypeInfo) *model.DishType {
	return &model.DishType{ID: info.DishTypeID, DishTypeName: info.DishTypeName, MasterType: info.MasterTypeID}
}

func ConvertToDishInfoList(daoList []*model.Dish, dishTypeMap map[uint32]*model.DishType) []*dto.DishInfo {
	retList := make([]*dto.DishInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, ConvertToDishInfo(dao, dishTypeMap))
	}
	return retList
}

func ConvertToDishInfo(dao *model.Dish, dishTypeMap map[uint32]*model.DishType) *dto.DishInfo {
	dishInfo := &dto.DishInfo{DishID: dao.ID, DishName: dao.DishName, Material: dao.Material, Price: dao.Price}
	typeInfo, ok := dishTypeMap[dao.DishType]
	if ok {
		dishInfo.DishTypeID = typeInfo.ID
		dishInfo.DishTypeName = typeInfo.DishTypeName
		if typeInfo.MasterType != 0 {
			masterType, ok := dishTypeMap[typeInfo.MasterType]
			if ok {
				dishInfo.MasterTypeName = masterType.DishTypeName
			}
		}
	}
	return dishInfo
}

func ConvertFromDishInfo(info *dto.DishInfo) *model.Dish {
	return &model.Dish{ID: info.DishID, DishName: info.DishName, DishType: info.DishTypeID,
		Price: info.Price, Material: info.Material}
}

func ConvertFromWeekMenuInfo(weekMenu *dto.WeekMenuDetailInfo) (*model.WeekMenu, error) {
	menuContent := make([]map[uint8][]uint32, 0)
	for _, menu := range weekMenu.MenuList {
		mealContent := convertToMenuContent(menu.MealList)
		menuContent = append(menuContent, mealContent)
	}

	content, err := json.Marshal(menuContent)
	if err != nil {
		logger.Warn(convertLogTag, "ConvertFromWeekMenu Failed|Err:%v", err)
		return nil, err
	}

	return &model.WeekMenu{ID: weekMenu.WeekMenuID, MenuTypeID: weekMenu.MenuType, MenuContent: string(content),
		MenuStartDate: time.Unix(weekMenu.MenuStartDate, 0)}, nil
}

func ConvertToWeekMenuList(daoList []*model.WeekMenu, dishList map[uint32]*model.Dish) []*dto.WeekMenuInfo {
	retList := make([]*dto.WeekMenuInfo, 0, len(daoList))
	for _, dao := range daoList {
		mealStrList, err := ConvertFromWeekMenuContent(dao.MenuContent, dishList)
		if err != nil {
			continue
		}
		retList = append(retList, &dto.WeekMenuInfo{WeekMenuID: dao.ID, MenuType: dao.MenuTypeID,
			MenuStartDate: dao.MenuStartDate.Unix(), MenuEndDate: dao.MenuStartDate.Add(time.Hour * 24 * 7).Unix(),
			MenuContent: mealStrList})
	}
	return retList
}

func ConvertToWeekMenuDetail(dao *model.WeekMenu, dishList map[uint32]*model.Dish,
	dishTypeMap map[uint32]*model.DishType) (*dto.WeekMenuDetailInfo, error) {
	menuList, err := ConvertDetailFromWeekMenuContent(dao.MenuStartDate, dao.MenuContent, dishList, dishTypeMap)
	if err != nil {
		logger.Warn(convertLogTag, "ConvertDetailFromWeekMenuContent Failed|Err:%v", err)
		return nil, err
	}

	return &dto.WeekMenuDetailInfo{
		WeekMenuID:    dao.ID,
		MenuType:      dao.MenuTypeID,
		MenuStartDate: dao.MenuStartDate.Unix(),
		MenuEndDate:   dao.MenuStartDate.Add(time.Hour * 24 * 7).Unix(),
		MenuList:      menuList,
	}, nil
}

func ConvertFromWeekMenuContent(content string, dishMap map[uint32]*model.Dish) ([]string, error) {
	contentMap := make([]map[uint8][]uint32, 0)
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		logger.Warn(convertLogTag, "ConvertFromMenuContent Failed|Err:%v", err)
		return nil, err
	}

	mealStrList := make([]string, 0)
	for _, dayMenu := range contentMap {
		for _, dishContent := range dayMenu {
			mealStr := ""
			for _, dishID := range dishContent {
				mealStr += dishMap[dishID].DishName + ","
			}
			mealStrList = append(mealStrList, mealStr[:len(mealStr)-1])
		}
	}

	return mealStrList, nil
}

func ConvertDetailFromWeekMenuContent(startDate time.Time, content string, dishMap map[uint32]*model.Dish,
	dishTypeMap map[uint32]*model.DishType) ([]*dto.MenuInfo, error) {
	contentMap := make([]map[uint8][]uint32, 0)
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		logger.Warn(convertLogTag, "ConvertFromMenuContent Failed|Err:%v", err)
		return nil, err
	}

	menuList := make([]*dto.MenuInfo, 0)
	for _, dayMenu := range contentMap {
		mealList, err := ConvertFromMenuContent2(dayMenu, dishMap, dishTypeMap)
		if err != nil {
			logger.Warn(convertLogTag, "ConvertDetailFromWeekMenuContent Failed|Err:%v", err)
			return nil, err
		}
		menuInfo := &dto.MenuInfo{MenuDate: startDate.Unix(), MealList: mealList}
		startDate = startDate.Add(time.Hour * 24)
		menuList = append(menuList, menuInfo)
	}

	return menuList, nil
}

func convertToMenuContent(mealList []*dto.MealInfo) map[uint8][]uint32 {
	contentMap := make(map[uint8][]uint32, len(mealList))
	for _, mealInfo := range mealList {
		contentMap[mealInfo.MealType] = make([]uint32, 0)
		for _, dish := range mealInfo.DishList {
			contentMap[mealInfo.MealType] = append(contentMap[mealInfo.MealType], dish.DishID)
		}
	}
	return contentMap
}

func ConvertFromMenuContent2(mealMap map[uint8][]uint32, dishMap map[uint32]*model.Dish,
	dishTypeMap map[uint32]*model.DishType) ([]*dto.MealInfo, error) {
	mealList := make([]*dto.MealInfo, 0)
	for mealType, dishContent := range mealMap {
		mealInfo := &dto.MealInfo{
			MealName: enum.GetMealName(mealType),
			MealType: mealType,
			DishList: make([]*dto.DishInfo, 0),
		}
		for _, dishID := range dishContent {
			dishInfo := ConvertToDishInfo(dishMap[dishID], dishTypeMap)
			mealInfo.DishList = append(mealInfo.DishList, dishInfo)
		}
		mealList = append(mealList, mealInfo)
	}
	return mealList, nil
}

func ConvertFromStoreTypeInfo(info *dto.StoreTypeInfo) *model.StorehouseType {
	return &model.StorehouseType{ID: info.StoreTypeID, StoreTypeName: info.StoreTypeName}
}

func ConvertToStoreTypeInfoList(daoList []*model.StorehouseType) []*dto.StoreTypeInfo {
	retList := make([]*dto.StoreTypeInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.StoreTypeInfo{StoreTypeID: dao.ID, StoreTypeName: dao.StoreTypeName})
	}
	return retList
}

func ConvertFromGoodsTypeInfo(info *dto.GoodsTypeInfo) *model.GoodsType {
	return &model.GoodsType{ID: info.GoodsTypeID, GoodsTypeName: info.GoodsTypeName, Discount: info.Discount}
}

func ConvertToGoodsTypeInfoList(daoList []*model.GoodsType) []*dto.GoodsTypeInfo {
	retList := make([]*dto.GoodsTypeInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.GoodsTypeInfo{GoodsTypeID: dao.ID, GoodsTypeName: dao.GoodsTypeName, Discount: dao.Discount})
	}
	return retList
}

func ConvertFromGoodsInfo(info *dto.GoodsInfo, picture string) *model.Goods {
	if picture == "" {
		picture = info.Picture
	}
	return &model.Goods{ID: info.GoodsID, Name: info.GoodsName, GoodsTypeID: info.GoodsType, StoreTypeID: info.StoreType,
		Picture: picture, BatchSize: info.BatchSize, BatchUnit: info.BatchUnit, Price: info.Price, Quantity: info.Quantity}
}

func ConvertToGoodsInfoList(daoList []*model.Goods) []*dto.GoodsInfo {
	retList := make([]*dto.GoodsInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.GoodsInfo{GoodsID: dao.ID, GoodsName: dao.Name, GoodsType: dao.GoodsTypeID,
			StoreType: dao.StoreTypeID, Picture: dao.Picture, BatchSize: dao.BatchSize, BatchUnit: dao.BatchUnit,
			Price: dao.Price, Quantity: dao.Quantity})
	}
	return retList
}
