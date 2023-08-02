package server

import (
	"fmt"

	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	IndexID           = "ID"
	IndexMenuDate     = "MenuDate"
	IndexMenuTypeName = "MenuTypeName"

	TableGeneratorLogTag = "TableGenerator"
)

func GenerateStaffMenuListTableHead() []*dto.TableColumnInfo {
	head := make([]*dto.TableColumnInfo, 0)
	head = append(head, &dto.TableColumnInfo{Name: "菜单ID", DataIndex: IndexID, Hide: false})
	head = append(head, &dto.TableColumnInfo{Name: "菜单日期", DataIndex: IndexMenuDate, Hide: false})
	for i := enum.MealUnknown + 1; i < enum.MealALL; i++ {
		head = append(head, &dto.TableColumnInfo{Name: enum.GetMealName(i), DataIndex: enum.GetMealKey(i), Hide: false})
	}
	return head
}

func GenerateStaffMenuListTableData(menuList []*model.Menu, dishIDMap map[uint32]*model.Dish) []*dto.TableRowInfo {
	dataList := make([]*dto.TableRowInfo, 0)
	for _, menu := range menuList {
		row := &dto.TableRowInfo{DataMap: make(map[string]*dto.TableRowColumnInfo)}
		row.DataMap[IndexID] = &dto.TableRowColumnInfo{Value: fmt.Sprintf("%v", menu.ID)}
		row.DataMap[IndexMenuDate] = &dto.TableRowColumnInfo{Value: menu.MenuDate.Format("2006-02-01")}
		dishMap, err := convertFromMenuContent(menu.MenuContent)
		if err != nil {
			logger.Warn(TableGeneratorLogTag, "convertFromMenuContent Failed|Err:%v", err)
			continue
		}
		for mealType, dishList := range dishMap {
			dishStr := ""
			for _, dishID := range dishList {
				dishStr += "," + dishIDMap[dishID].DishName
			}
			if len(dishList) == 0 {
				dishStr = ","
			}
			row.DataMap[enum.GetMealKey(mealType)] = &dto.TableRowColumnInfo{Value: dishStr[1:]}
		}
		dataList = append(dataList, row)
	}
	return dataList
}

func GenerateStaffDetailTableHead() []*dto.TableColumnInfo {
	head := make([]*dto.TableColumnInfo, 0)
	for i := enum.MealUnknown + 1; i < enum.MealALL; i++ {
		head = append(head, &dto.TableColumnInfo{DataIndex: enum.GetMealKey(i), Hide: true})
		head = append(head, &dto.TableColumnInfo{DataIndex: enum.GetMealKey(i) + "DishType", Hide: true})
		head = append(head, &dto.TableColumnInfo{DataIndex: enum.GetMealKey(i) + "Dish", Hide: true})
	}
	return head
}

func GenerateStaffDetailTableData(menu *model.Menu, dishIDMap map[uint32]*model.Dish,
	dishTypeMap map[uint32]*model.DishType) []*dto.TableRowInfo {
	dishMap, err := convertFromMenuContent(menu.MenuContent)
	if err != nil {
		logger.Warn(TableGeneratorLogTag, "convertFromMenuContent Failed|Err:%v", err)
		return nil
	}
	logger.Debug(TableGeneratorLogTag, "StaffDetail|DishMap:%#v", dishMap)

	maxLen := 0
	mealByType := make(map[uint8]map[uint32][]*model.Dish)
	mealDishNumber := make(map[uint8]int)
	for mealType, dishList := range dishMap {
		dishByType := make(map[uint32][]*model.Dish)
		mealDishNumber[mealType] = len(dishList)
		if maxLen < len(dishList) {
			maxLen = len(dishList)
		}

		for _, dishID := range dishList {
			dish := dishIDMap[dishID]
			if _, ok := dishByType[dish.DishType]; ok == false {
				dishByType[dish.DishType] = make([]*model.Dish, 0)
			}
			dishByType[dish.DishType] = append(dishByType[dish.DishType], dish)
		}
		mealByType[mealType] = dishByType
	}

	rows := make([]*dto.TableRowInfo, maxLen)
	for index := range rows {
		rows[index] = &dto.TableRowInfo{DataMap: make(map[string]*dto.TableRowColumnInfo)}
	}

	logger.Debug(TableGeneratorLogTag, "StaffDetail|mealByType:%#v", mealByType)
	for mealType, dishByType := range mealByType {
		curIndex := 0
		extraRow := maxLen % len(dishMap[mealType])
		repeatTimes := maxLen / len(dishMap[mealType])

		for dishType, dishList := range dishByType {
			logger.Debug(TableGeneratorLogTag, "StaffDetail|dishType:%v|Len:%v", dishByType, len(dishList))
			for _, dish := range dishList {
				rows[curIndex].DataMap[enum.GetMealKey(mealType)] = &dto.TableRowColumnInfo{
					Value: fmt.Sprintf("%v(%v)", enum.GetMealName(mealType), mealDishNumber[mealType])}
				rows[curIndex].DataMap[enum.GetMealKey(mealType)+"DishType"] = &dto.TableRowColumnInfo{
					Value: fmt.Sprintf("%v(%v)", dishTypeMap[dishType].DishTypeName, len(dishList))}
				rows[curIndex].DataMap[enum.GetMealKey(mealType)+"Dish"] = &dto.TableRowColumnInfo{
					Value: dish.DishName}
				curIndex++
				extraTimes := 0
				if extraRow > 0 {
					extraTimes = 1
					extraRow--
				}
				for i := 1; i < repeatTimes+extraTimes; i++ {
					rows[curIndex].DataMap[enum.GetMealKey(mealType)] =
						rows[curIndex-1].DataMap[enum.GetMealKey(mealType)]
					rows[curIndex].DataMap[enum.GetMealKey(mealType)+"DishType"] =
						rows[curIndex-1].DataMap[enum.GetMealKey(mealType)+"DishType"]
					rows[curIndex].DataMap[enum.GetMealKey(mealType)+"Dish"] =
						rows[curIndex-1].DataMap[enum.GetMealKey(mealType)+"Dish"]
					curIndex++
				}
			}
		}
	}

	return rows
}

func GenerateMenuTypeListTableHead() []*dto.TableColumnInfo {
	head := make([]*dto.TableColumnInfo, 0)
	head = append(head, &dto.TableColumnInfo{Name: "菜单类型ID", DataIndex: IndexID, Hide: false})
	head = append(head, &dto.TableColumnInfo{Name: "菜单类型名称", DataIndex: IndexMenuTypeName, Hide: false})
	for i := enum.MealUnknown + 1; i < enum.MealALL; i++ {
		head = append(head, &dto.TableColumnInfo{Name: enum.GetMealName(i), DataIndex: enum.GetMealKey(i), Hide: false})
	}

	return head
}

func GenerateMenuTypeListTableData(menuTypeList []*model.MenuType, dishTypeMap map[uint32]*model.DishType) []*dto.TableRowInfo {
	dataList := make([]*dto.TableRowInfo, 0)
	for _, menuType := range menuTypeList {
		row := &dto.TableRowInfo{DataMap: make(map[string]*dto.TableRowColumnInfo)}
		row.DataMap[IndexID] = &dto.TableRowColumnInfo{Value: fmt.Sprintf("%v", menuType.ID)}
		row.DataMap[IndexMenuTypeName] = &dto.TableRowColumnInfo{Value: menuType.MenuTypeName}
		dataList = append(dataList, row)

		menuContentMap := menuType.ToMenuConfig()
		if menuContentMap == nil {
			continue
		}

		for mealType, mealContent := range menuContentMap {
			content := ""
			for dishType, dishNumber := range mealContent {
				content += fmt.Sprintf(",%v:%v", dishTypeMap[dishType].DishTypeName, dishNumber)
			}
			if content == "" {
				content = ","
			}
			row.DataMap[enum.GetMealKey(mealType)] = &dto.TableRowColumnInfo{Value: content[1:]}
		}
	}
	return dataList
}

func GenerateMenuTypeDetailTableHead(menuType *model.MenuType, dishTypeMap map[uint32]*model.DishType) []*dto.TableColumnInfo {
	head := make([]*dto.TableColumnInfo, 0)
	menuContentMap := menuType.ToMenuConfig()
	if menuContentMap == nil {
		return nil
	}

	for mealType, mealContent := range menuContentMap {
		mealHead := &dto.TableColumnInfo{Name: enum.GetMealName(mealType), DataIndex: enum.GetMealKey(mealType),
			Hide: false, Children: make([]*dto.TableColumnInfo, 0)}
		for dishType := range mealContent {
			dishHead := &dto.TableColumnInfo{Name: dishTypeMap[dishType].DishTypeName,
				DataIndex: enum.GetMealKey(mealType) + fmt.Sprintf("%v", dishTypeMap[dishType].ID), Hide: false}
			mealHead.Children = append(mealHead.Children, dishHead)
		}
		head = append(head, mealHead)
	}

	return head
}

func GenerateMenuTypeDetailTableData(menuType *model.MenuType, dishTypeMap map[uint32]*model.DishType) []*dto.TableRowInfo {
	dataList := make([]*dto.TableRowInfo, 0)
	menuContentMap := menuType.ToMenuConfig()
	if menuContentMap == nil {
		return nil
	}

	row := &dto.TableRowInfo{DataMap: make(map[string]*dto.TableRowColumnInfo)}
	for mealType, mealContent := range menuContentMap {
		for dishType, dishNumber := range mealContent {
			row.DataMap[enum.GetMealKey(mealType)+fmt.Sprintf("%v", dishTypeMap[dishType].ID)] =
				&dto.TableRowColumnInfo{Value: fmt.Sprintf("%v", dishNumber)}
		}
	}
	dataList = append(dataList, row)
	return dataList
}
