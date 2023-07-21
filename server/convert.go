package server

import (
	"encoding/json"
	"time"

	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

func ConvertToDishTypeInfoList(daoList []*model.DishType) []*dto.DishTypeInfo {
	retList := make([]*dto.DishTypeInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.DishTypeInfo{DishTypeID: dao.ID, DishTypeName: dao.DishTypeName})
	}
	return retList
}

func ConvertFromDishTypeInfo(info *dto.DishTypeInfo) *model.DishType {
	return &model.DishType{ID: info.DishTypeID, DishTypeName: info.DishTypeName}
}

func ConvertToDishInfoList(daoList []*model.Dish) []*dto.DishInfo {
	retList := make([]*dto.DishInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, ConvertToDishInfo(dao))
	}
	return retList
}
func ConvertToDishInfo(dao *model.Dish) *dto.DishInfo {
	return &dto.DishInfo{DishID: dao.ID, DishName: dao.DishName,
		DishType: dao.DishType, Material: dao.Material, Price: dao.Price}
}

func ConvertFromDishInfo(info *dto.DishInfo) *model.Dish {
	return &model.Dish{ID: info.DishID, DishName: info.DishName, DishType: info.DishType,
		Price: info.Price, Material: info.Material}
}

func ConvertToMenuInfoList(daoList []*model.Menu, dishList map[uint32]*model.Dish) []*dto.MenuInfo {
	retList := make([]*dto.MenuInfo, 0, len(daoList))
	for _, dao := range daoList {
		mealList, err := ConvertFromMenuContent(dao.MenuContent, dishList)
		if err != nil {
			continue
		}
		retList = append(retList, &dto.MenuInfo{MenuID: dao.ID, MenuType: dao.MenuTypeID, MealList: mealList,
			MenuDate: dao.MenuDate.Unix()})
	}
	return retList
}

func ConvertFromMenuInfo(menuInfo *dto.MenuInfo) (*model.Menu, error) {
	content, err := ConvertToMenuContent(menuInfo.MealList)
	if err != nil {
		return nil, err
	}
	return &model.Menu{ID: menuInfo.MenuID, MenuTypeID: menuInfo.MenuType, MenuContent: content,
		MenuDate: time.Unix(menuInfo.MenuDate, 0)}, nil
}

func ConvertToMenuContent(mealList []*dto.MealInfo) (string, error) {
	contentMap := make(map[uint8][]uint32, len(mealList))
	for _, mealInfo := range mealList {
		contentMap[mealInfo.MealType] = make([]uint32, 0)
		for _, dish := range mealInfo.DishList {
			contentMap[mealInfo.MealType] = append(contentMap[mealInfo.MealType], dish.DishID)
		}
	}
	content, err := json.Marshal(contentMap)
	if err != nil {
		logger.Warn(dishServerLogTag, "ConvertToMenuContent Failed|Err:%v", err)
		return "", err
	}
	return string(content), nil
}

func ConvertFromMenuContent(content string, dishMap map[uint32]*model.Dish) ([]*dto.MealInfo, error) {
	contentMap := make(map[uint8][]uint32, 0)
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		logger.Warn(dishServerLogTag, "ConvertFromMenuContent Failed|Err:%v", err)
		return nil, err
	}

	mealList := make([]*dto.MealInfo, 0)
	for mealType, dishContent := range contentMap {
		mealInfo := &dto.MealInfo{
			MealName: enum.GetMealName(mealType),
			MealType: mealType,
			DishList: make([]*dto.DishInfo, 0),
		}
		for _, dishID := range dishContent {
			dishInfo := ConvertToDishInfo(dishMap[dishID])
			mealInfo.DishList = append(mealInfo.DishList, dishInfo)
		}
		mealList = append(mealList, mealInfo)
	}
	return mealList, nil
}

func ConvertToMenuTypeList(daoList []*model.MenuType) []*dto.MenuTypeInfo {
	retList := make([]*dto.MenuTypeInfo, 0, len(daoList))
	for _, dao := range daoList {
		menuConf, err := ConvertFromMenuTypeConfig(dao.MenuConfig)
		if err != nil {
			logger.Warn(dishServerLogTag, "ConvertFromMenuTypeConfig Failed|Err:%v", err)
			continue
		}
		retList = append(retList, &dto.MenuTypeInfo{MenuTypeID: dao.ID, Name: dao.MenuTypeName,
			MenuConfig: menuConf})
	}
	return retList
}

func ConvertFromMenuTypeInfo(typeInfo *dto.MenuTypeInfo) (*model.MenuType, error) {
	menuConf, err := ConvertToMenuTypeConfig(typeInfo.MenuConfig)
	if err != nil {
		logger.Warn(dishServerLogTag, "ConvertToMenuTypeConfig Failed|Err:%v", err)
		return nil, err
	}
	return &model.MenuType{ID: typeInfo.MenuTypeID, MenuTypeName: typeInfo.Name, MenuConfig: menuConf}, nil
}

func ConvertToMenuTypeConfig(menuConfList []*dto.MenuConfigInfo) (string, error) {
	menuConfig := make(map[uint8]map[uint32]int32, 0)
	for _, confInfo := range menuConfList {
		menuConfig[confInfo.MealType] = confInfo.DishNumberMap
	}
	menuConf, err := json.Marshal(menuConfig)
	if err != nil {
		logger.Warn(dishServerLogTag, "ConvertToMenuTypeConfig Failed|Err:%v", err)
		return "", err
	}
	return string(menuConf), nil
}

func ConvertFromMenuTypeConfig(menuConf string) ([]*dto.MenuConfigInfo, error) {
	menuConfig := make(map[uint8]map[uint32]int32, 0)
	err := json.Unmarshal([]byte(menuConf), &menuConfig)
	if err != nil {
		logger.Warn(dishServerLogTag, "ConvertFromMenuTypeConfig Failed|Err:%v", err)
		return nil, err
	}

	retList := make([]*dto.MenuConfigInfo, 0)
	for mealType, dishNumberMap := range menuConfig {
		confInfo := &dto.MenuConfigInfo{
			MealName:      enum.GetMealName(mealType),
			MealType:      mealType,
			DishNumberMap: dishNumberMap,
		}
		retList = append(retList, confInfo)
	}
	return retList, nil
}
