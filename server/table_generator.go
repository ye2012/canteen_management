package server

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	IndexID           = "ID"
	IndexMenuDate     = "MenuDate"
	IndexMenuTypeName = "MenuTypeName"
	IndexDish         = "Dish"
	IndexDishType     = "DishType"

	IndexDelimiter = "_"

	TableGeneratorLogTag = "TableGenerator"

	DayRowFixed = 4
)

var (
	weekDays = []string{"星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期日"}
)

func GenerateDishIndex(mealType enum.MealType) string {
	return enum.GetMealKey(mealType) + IndexDelimiter + IndexDish
}

func GenerateDishTypeIndex(mealType enum.MealType) string {
	return enum.GetMealKey(mealType) + IndexDelimiter + IndexDishType
}

func GenerateDayMealIndex(mealType enum.MealType, dayIndex int) string {
	return enum.GetMealKey(mealType) + IndexDelimiter + fmt.Sprintf("%v", dayIndex)
}

func GenerateMealMainDishTypeIndex(mealType enum.MealType, mainTypeID uint32, index int) string {
	return fmt.Sprintf("%v%v%v%v%v", enum.GetMealKey(mealType), IndexDelimiter, mainTypeID, IndexDelimiter, index)
}

func GenerateWeekDayValue(date time.Time) string {
	return fmt.Sprintf("%v\n(%v月%v日)", weekDays[(date.Weekday()+6)%7], int(date.Month()), date.Day())
}

func GenerateStaffMenuListTableHead() []*dto.TableColumnInfo {
	head := make([]*dto.TableColumnInfo, 0)
	head = append(head, &dto.TableColumnInfo{Name: "菜单ID", DataIndex: IndexID, Hide: false, MergeRow: false})
	head = append(head, &dto.TableColumnInfo{Name: "菜单日期", DataIndex: IndexMenuDate, Hide: false, MergeRow: false})
	for i := enum.MealUnknown + 1; i < enum.MealALL; i++ {
		head = append(head, &dto.TableColumnInfo{Name: enum.GetMealName(i), DataIndex: enum.GetMealKey(i), Hide: false,
			MergeRow: false})
	}
	return head
}

func GenerateStaffMenuListTableData(menuList []*model.Menu, dishIDMap map[uint32]*model.Dish) []dto.TableRowInfo {
	dataList := make([]dto.TableRowInfo, 0)
	for _, menu := range menuList {
		row := make(map[string]*dto.TableRowColumnInfo)
		row[IndexID] = &dto.TableRowColumnInfo{ID: menu.ID, Value: fmt.Sprintf("%v", menu.ID)}
		row[IndexMenuDate] = &dto.TableRowColumnInfo{Value: menu.MenuDate.Format("2006-01-02")}
		dishMap := menu.ToMenuContent()
		if dishMap == nil {
			logger.Warn(TableGeneratorLogTag, "StaffMenuList ToMenuContent Failed")
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
			row[enum.GetMealKey(mealType)] = &dto.TableRowColumnInfo{Value: dishStr[1:]}
		}
		dataList = append(dataList, row)
	}
	return dataList
}

func GenerateStaffDetailTableHead() []*dto.TableColumnInfo {
	head := make([]*dto.TableColumnInfo, 0)
	for i := enum.MealUnknown + 1; i < enum.MealALL; i++ {
		head = append(head, &dto.TableColumnInfo{DataIndex: enum.GetMealKey(i), Hide: true, MergeRow: true})
		head = append(head, &dto.TableColumnInfo{DataIndex: GenerateDishTypeIndex(i), Hide: true, MergeRow: true})
		head = append(head, &dto.TableColumnInfo{DataIndex: GenerateDishIndex(i), Hide: true, MergeRow: true})
	}
	return head
}

func GenerateStaffDetailTableData(menu *model.Menu, dishIDMap map[uint32]*model.Dish,
	dishTypeMap map[uint32]*model.DishType) []dto.TableRowInfo {
	dishMap := menu.ToMenuContent()
	if dishMap == nil {
		logger.Warn(TableGeneratorLogTag, "ToStaffMenuContent Failed")
		return nil
	}

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

	rows := make([]dto.TableRowInfo, maxLen)
	for index := range rows {
		rows[index] = make(map[string]*dto.TableRowColumnInfo)
	}

	logger.Debug(TableGeneratorLogTag, "StaffDetail|mealByType:%#v", mealByType)
	for mealType, dishByType := range mealByType {
		curIndex := 0
		extraRow := maxLen % len(dishMap[mealType])
		repeatTimes := maxLen / len(dishMap[mealType])

		for dishType, dishList := range dishByType {
			logger.Debug(TableGeneratorLogTag, "StaffDetail|dishType:%v|Len:%v", dishByType, len(dishList))
			for _, dish := range dishList {
				rows[curIndex][enum.GetMealKey(mealType)] = &dto.TableRowColumnInfo{
					Value: fmt.Sprintf("%v(%v)", enum.GetMealName(mealType), mealDishNumber[mealType])}
				rows[curIndex][GenerateDishTypeIndex(mealType)] = &dto.TableRowColumnInfo{
					Value: fmt.Sprintf("%v(%v)", dishTypeMap[dishType].DishTypeName, len(dishList))}
				rows[curIndex][GenerateDishIndex(mealType)] = &dto.TableRowColumnInfo{
					ID: dish.ID, Value: dish.DishName}
				curIndex++
				extraTimes := 0
				if extraRow > 0 {
					extraTimes = 1
					extraRow--
				}
				for i := 1; i < repeatTimes+extraTimes; i++ {
					rows[curIndex][enum.GetMealKey(mealType)] =
						rows[curIndex-1][enum.GetMealKey(mealType)]
					rows[curIndex][GenerateDishTypeIndex(mealType)] =
						rows[curIndex-1][GenerateDishTypeIndex(mealType)]
					rows[curIndex][GenerateDishIndex(mealType)] =
						rows[curIndex-1][GenerateDishIndex(mealType)]
					curIndex++
				}
			}
		}
	}
	return rows
}

func ParseStaffMenuDetailData(rowList []dto.ModifyTableRowInfo, menuID, menuTypeID uint32, menuDate int64) *model.Menu {
	menu := &model.Menu{ID: menuID, MenuTypeID: menuTypeID, MenuDate: time.Unix(menuDate, 0)}
	menuConf := make(map[uint8][]uint32)
	mealDishMap := make(map[uint8]map[uint32]bool)
	for _, row := range rowList {
		for key, value := range row {
			keys := strings.Split(key, IndexDelimiter)
			if len(keys) <= 1 {
				continue
			}
			mealType := enum.GetMealType(keys[0])
			if _, ok := mealDishMap[mealType]; ok == false {
				mealDishMap[mealType] = make(map[uint32]bool)
			}
			if mealType == enum.MealUnknown {
				logger.Warn(TableGeneratorLogTag, "StaffMenuDetail Unknown Meal|Key:%v|Value:%#v", key, value)
				continue
			}
			if keys[1] == IndexDish {
				if _, ok := menuConf[mealType]; ok == false {
					menuConf[mealType] = make([]uint32, 0)
				}
				dishID := uint32(value.(float64))
				if _, ok := mealDishMap[mealType][dishID]; ok {
					continue
				}
				mealDishMap[mealType][dishID] = true
				menuConf[mealType] = append(menuConf[mealType], dishID)
			}
		}
	}
	menu.FromMenuConfig(menuConf)
	return menu
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

func GenerateMenuTypeListTableData(menuTypeList []*model.MenuType, dishTypeMap map[uint32]*model.DishType) []dto.TableRowInfo {
	dataList := make([]dto.TableRowInfo, 0)
	for _, menuType := range menuTypeList {
		row := make(map[string]*dto.TableRowColumnInfo)
		row[IndexID] = &dto.TableRowColumnInfo{ID: menuType.ID, Value: fmt.Sprintf("%v", menuType.ID)}
		row[IndexMenuTypeName] = &dto.TableRowColumnInfo{Value: menuType.MenuTypeName}
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
			row[enum.GetMealKey(mealType)] = &dto.TableRowColumnInfo{Value: content[1:]}
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

	for mealType := enum.MealUnknown + 1; mealType < enum.MealALL; mealType++ {
		mealContent, ok := menuContentMap[mealType]
		if ok == false {
			continue
		}
		mealHead := &dto.TableColumnInfo{Name: enum.GetMealName(mealType), DataIndex: enum.GetMealKey(mealType),
			Hide: false, Children: make([]*dto.TableColumnInfo, 0)}
		for dishType := range mealContent {
			dishHead := &dto.TableColumnInfo{Name: dishTypeMap[dishType].DishTypeName, MergeRow: true,
				DataIndex: enum.GetMealKey(mealType) + fmt.Sprintf("_%v", dishTypeMap[dishType].ID), Hide: false}
			mealHead.Children = append(mealHead.Children, dishHead)
		}
		head = append(head, mealHead)
	}

	return head
}

func GenerateMenuTypeDetailTableData(menuType *model.MenuType, dishTypeMap map[uint32]*model.DishType) []dto.TableRowInfo {
	dataList := make([]dto.TableRowInfo, 0)
	menuContentMap := menuType.ToMenuConfig()
	if menuContentMap == nil {
		return nil
	}

	row := make(map[string]*dto.TableRowColumnInfo)
	row[IndexID] = &dto.TableRowColumnInfo{ID: menuType.ID, Value: fmt.Sprintf("%v", menuType.ID)}
	row[IndexMenuTypeName] = &dto.TableRowColumnInfo{Value: menuType.MenuTypeName}
	for mealType, mealContent := range menuContentMap {
		for dishType, dishNumber := range mealContent {
			row[enum.GetMealKey(mealType)+fmt.Sprintf("_%v", dishTypeMap[dishType].ID)] =
				&dto.TableRowColumnInfo{Value: fmt.Sprintf("%v", dishNumber)}
		}
	}
	dataList = append(dataList, row)
	return dataList
}

func ParseMenuTypeDetailData(row dto.ModifyTableRowInfo, menuTypeID uint32, menuTypeName string) *model.MenuType {
	menuType := &model.MenuType{ID: menuTypeID, MenuTypeName: menuTypeName}
	menuConf := make(map[uint8]map[uint32]int32)
	for key, dishNumber := range row {
		keys := strings.Split(key, IndexDelimiter)
		if len(keys) <= 1 {
			continue
		}
		mealType := enum.GetMealType(keys[0])
		if mealType == enum.MealUnknown {
			logger.Warn(TableGeneratorLogTag, "MenuTypeDetail Unknown Meal|Key:%v|Number:%#v", key, dishNumber)
			continue
		}
		if _, ok := menuConf[mealType]; ok == false {
			menuConf[mealType] = make(map[uint32]int32)
		}
		dishType, _ := strconv.ParseInt(keys[1], 10, 32)
		menuConf[mealType][uint32(dishType)] = int32(dishNumber.(float64))
	}
	menuType.FromMenuConfig(menuConf)
	return menuType
}

func GenerateWeekMenuListTableHead() []*dto.TableColumnInfo {
	head := make([]*dto.TableColumnInfo, 0)
	head = append(head, &dto.TableColumnInfo{Name: "菜单ID", DataIndex: IndexID, Hide: false})
	head = append(head, &dto.TableColumnInfo{Name: "菜单日期", DataIndex: IndexMenuDate, Hide: false})
	for index, dayName := range weekDays {
		for i := enum.MealUnknown + 1; i < enum.MealALL; i++ {
			head = append(head, &dto.TableColumnInfo{Name: enum.GetMealName(i) + dayName,
				DataIndex: GenerateDayMealIndex(i, index), Hide: false})
		}
	}

	return head
}

func GenerateWeekMenuListTableData(menuList []*model.WeekMenu, dishIDMap map[uint32]*model.Dish) []dto.TableRowInfo {
	dataList := make([]dto.TableRowInfo, 0)
	for _, menu := range menuList {
		row := make(map[string]*dto.TableRowColumnInfo)
		row[IndexID] = &dto.TableRowColumnInfo{ID: menu.ID, Value: fmt.Sprintf("%v", menu.ID)}
		startDate := menu.MenuStartDate.Format("01-02")
		endDate := menu.MenuStartDate.Add(time.Hour * 24 * 7).Format("01-02")
		row[IndexMenuDate] = &dto.TableRowColumnInfo{Value: fmt.Sprintf("%v~%v", startDate, endDate)}
		dataList = append(dataList, row)

		menuConf := menu.ToWeekMenuConfig()
		if menuConf == nil {
			continue
		}
		for dayIndex, dayConf := range menuConf {
			for mealType, dishList := range dayConf {
				mealStr := ""
				for _, dishID := range dishList {
					mealStr += dishIDMap[dishID].DishName + ","
				}
				if len(dishList) == 0 {
					mealStr = ","
				}
				row[GenerateDayMealIndex(mealType, dayIndex%len(weekDays))] =
					&dto.TableRowColumnInfo{Value: mealStr[:len(mealStr)-1]}
			}
		}
	}
	return dataList
}

func GenerateWeekMenuDetailTable(menu *model.WeekMenu, dishIDMap map[uint32]*model.Dish,
	dishTypeMap map[uint32]*model.DishType) ([]*dto.TableColumnInfo, []dto.TableRowInfo) {
	head := make([]*dto.TableColumnInfo, 0)
	dataList := make([]dto.TableRowInfo, 0)
	head = append(head, &dto.TableColumnInfo{Name: "日期", DataIndex: IndexMenuDate, Hide: false, MergeRow: true})

	mainDishType := make([]*model.DishType, 0)
	for _, dishType := range dishTypeMap {
		if dishType.MasterType == 0 {
			mainDishType = append(mainDishType, dishType)
		}
	}

	typeDishList := make([]map[uint8]map[uint32]*list.List, 0)

	columnNumber := make(map[uint8]map[uint32]int32)
	menuConf := menu.ToWeekMenuConfig()
	for _, dayMenu := range menuConf {
		dayTypeDish := make(map[uint8]map[uint32]*list.List)
		for mealType, dishList := range dayMenu {
			if _, ok := columnNumber[mealType]; !ok {
				columnNumber[mealType] = make(map[uint32]int32)
			}
			dayTypeDish[mealType] = make(map[uint32]*list.List)
			typeNumber := make(map[uint32]int32)
			for _, dishID := range dishList {
				dish := dishIDMap[dishID]
				dishType := dish.DishType
				mainType := dishTypeMap[dishType].MasterType
				if mainType == 0 {
					mainType = dishType
				}
				if _, ok := typeNumber[mainType]; !ok {
					typeNumber[mainType] = 0
				}
				typeNumber[mainType]++
				if _, ok := dayTypeDish[mealType][mainType]; !ok {
					dayTypeDish[mealType][mainType] = list.New()
				}
				dayTypeDish[mealType][mainType].PushBack(dish)
			}

			for mainType, dishNumber := range typeNumber {
				if _, ok := columnNumber[mealType][mainType]; !ok {
					columnNumber[mealType][mainType] = 0
				}
				number := dishNumber / DayRowFixed
				if dishNumber%DayRowFixed > 0 {
					number++
				}
				if columnNumber[mealType][mainType] < number {
					columnNumber[mealType][mainType] = number
				}
			}
		}
		typeDishList = append(typeDishList, dayTypeDish)
	}

	for mealType := enum.MealUnknown + 1; mealType < enum.MealALL; mealType++ {
		mealColumnNumber, ok := columnNumber[mealType]
		if ok == false {
			continue
		}
		mealColumn := &dto.TableColumnInfo{Name: enum.GetMealName(mealType),
			DataIndex: enum.GetMealKey(mealType), Hide: false}
		children := make([]*dto.TableColumnInfo, 0)
		for _, mainType := range mainDishType {
			number, ok := mealColumnNumber[mainType.ID]
			if ok == false {
				continue
			}
			for i := 0; i < int(number); i++ {
				children = append(children, &dto.TableColumnInfo{Name: mainType.DishTypeName, MergeRow: false,
					DataIndex: GenerateMealMainDishTypeIndex(mealType, mainType.ID, i), Hide: false})
			}
		}
		mealColumn.Children = children
		head = append(head, mealColumn)
	}

	startDate := menu.MenuStartDate
	for _, dayDishMap := range typeDishList {
		for rowIndex := 0; rowIndex < DayRowFixed; rowIndex++ {
			row := dto.TableRowInfo{}
			row[IndexMenuDate] = &dto.TableRowColumnInfo{Value: GenerateWeekDayValue(startDate)}
			for mealType, dishTypeList := range dayDishMap {
				if _, ok := columnNumber[mealType]; ok == false {
					continue
				}
				for mainType, dishList := range dishTypeList {
					number := columnNumber[mealType][mainType]
					for column := 0; column < int(number); column++ {
						if dishList.Len() == 0 {
							row[GenerateMealMainDishTypeIndex(mealType, mainType, column)] =
								&dto.TableRowColumnInfo{ID: 0, Value: ""}
							continue
						}
						dishNode := dishList.Front()
						dishList.Remove(dishNode)
						dish := dishNode.Value.(*model.Dish)
						row[GenerateMealMainDishTypeIndex(mealType, mainType, column)] =
							&dto.TableRowColumnInfo{ID: dish.ID, Value: dish.DishName}
					}
				}
			}
			dataList = append(dataList, row)
		}
		startDate = startDate.Add(time.Hour * 24)
	}

	return head, dataList
}

func ParseWeekMenuDetail(rowList []dto.ModifyTableRowInfo, menuID, menuTypeID uint32, menuDate int64) *model.WeekMenu {
	weekMenu := &model.WeekMenu{ID: menuID, MenuTypeID: menuTypeID, MenuStartDate: time.Unix(menuDate, 0)}
	menuConf := make([]map[uint8][]uint32, 7)
	for rowIndex, row := range rowList {
		for key, value := range row {
			keys := strings.Split(key, IndexDelimiter)
			if len(keys) < 3 {
				continue
			}
			dateIndex := rowIndex / DayRowFixed
			mealType, _ := strconv.ParseInt(keys[0], 10, 8)
			if dateIndex > len(menuConf) {
				logger.Warn(TableGeneratorLogTag, "DateIndex Extent MenuConfLen|rowIndex:%v|ConfLen:%v",
					rowIndex, len(menuConf))
				continue
			}
			if _, ok := menuConf[dateIndex][uint8(mealType)]; ok == false {
				menuConf[dateIndex][uint8(mealType)] = make([]uint32, 0)
			}
			dishID := uint32(value.(float64))
			menuConf[dateIndex][uint8(mealType)] = append(menuConf[dateIndex][uint8(mealType)], dishID)
		}
	}
	weekMenu.FromWeekMenuConfig(menuConf)
	return weekMenu
}
