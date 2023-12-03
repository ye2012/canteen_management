package conv

import (
	"encoding/json"
	"github.com/canteen_management/enum"
	"time"

	"github.com/canteen_management/dto"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	dishConvertLogTag = "DishConvert"
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
	dishInfo := &dto.DishInfo{DishID: dao.ID, DishName: dao.DishName, Picture: dao.Picture,
		Material: dao.Material, Price: dao.Price}
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
		Picture: info.Picture, Price: info.Price, Material: info.Material}
}

func ConvertFromDishList(infoList []*dto.DishInfo) []*model.Dish {
	retList := make([]*model.Dish, 0, len(infoList))
	for _, dishInfo := range infoList {
		retList = append(retList, ConvertFromDishInfo(dishInfo))
	}
	return retList
}

func ConvertFromWeekMenuInfo(weekMenu *dto.WeekMenuDetailInfo) (*model.WeekMenu, error) {
	menuContent := make([]map[uint8][]uint32, 0)
	for _, menu := range weekMenu.MenuList {
		mealContent := convertToMenuContent(menu.MealList)
		menuContent = append(menuContent, mealContent)
	}

	content, err := json.Marshal(menuContent)
	if err != nil {
		logger.Warn(dishConvertLogTag, "ConvertFromWeekMenu Failed|Err:%v", err)
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
		logger.Warn(dishConvertLogTag, "ConvertDetailFromWeekMenuContent Failed|Err:%v", err)
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
		logger.Warn(dishConvertLogTag, "ConvertFromMenuContent Failed|Err:%v", err)
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
		logger.Warn(dishConvertLogTag, "ConvertFromMenuContent Failed|Err:%v", err)
		return nil, err
	}

	menuList := make([]*dto.MenuInfo, 0)
	for _, dayMenu := range contentMap {
		mealList, err := ConvertFromMenuContent2(dayMenu, dishMap, dishTypeMap)
		if err != nil {
			logger.Warn(dishConvertLogTag, "ConvertDetailFromWeekMenuContent Failed|Err:%v", err)
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
