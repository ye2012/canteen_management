package server

import (
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

func GenerateDayMealIndex(mealType enum.MealType, dayName string) string {
	return enum.GetMealKey(mealType) + IndexDelimiter + dayName
}

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
		row.DataMap[IndexID] = &dto.TableRowColumnInfo{ID: menu.ID, Value: fmt.Sprintf("%v", menu.ID)}
		row.DataMap[IndexMenuDate] = &dto.TableRowColumnInfo{Value: menu.MenuDate.Format("2006-01-02")}
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
		head = append(head, &dto.TableColumnInfo{DataIndex: GenerateDishTypeIndex(i), Hide: true})
		head = append(head, &dto.TableColumnInfo{DataIndex: GenerateDishIndex(i), Hide: true})
	}
	return head
}

func GenerateStaffDetailTableData(menu *model.Menu, dishIDMap map[uint32]*model.Dish,
	dishTypeMap map[uint32]*model.DishType) []*dto.TableRowInfo {
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
				rows[curIndex].DataMap[GenerateDishTypeIndex(mealType)] = &dto.TableRowColumnInfo{
					Value: fmt.Sprintf("%v(%v)", dishTypeMap[dishType].DishTypeName, len(dishList))}
				rows[curIndex].DataMap[GenerateDishIndex(mealType)] = &dto.TableRowColumnInfo{
					ID: dish.ID, Value: dish.DishName}
				curIndex++
				extraTimes := 0
				if extraRow > 0 {
					extraTimes = 1
					extraRow--
				}
				for i := 1; i < repeatTimes+extraTimes; i++ {
					rows[curIndex].DataMap[enum.GetMealKey(mealType)] =
						rows[curIndex-1].DataMap[enum.GetMealKey(mealType)]
					rows[curIndex].DataMap[GenerateDishTypeIndex(mealType)] =
						rows[curIndex-1].DataMap[GenerateDishTypeIndex(mealType)]
					rows[curIndex].DataMap[GenerateDishIndex(mealType)] =
						rows[curIndex-1].DataMap[GenerateDishIndex(mealType)]
					curIndex++
				}
			}
		}
	}
	return rows
}

func ParseStaffMenuDetailData(rowList []*dto.TableRowInfo, menuID uint32) *model.Menu {
	menu := &model.Menu{ID: menuID}
	menuConf := make(map[uint8][]uint32)
	for _, row := range rowList {
		for key, value := range row.DataMap {
			keys := strings.Split(key, IndexDelimiter)
			if len(keys) <= 1 {
				continue
			}
			mealType := enum.GetMealType(keys[0])
			if mealType == enum.MealUnknown {
				logger.Warn(TableGeneratorLogTag, "StaffMenuDetail Unknown Meal|Key:%v|Value:%#v", key, value)
				continue
			}
			if keys[1] == IndexDish && value.ID > 0 {
				if _, ok := menuConf[mealType]; ok == false {
					menuConf[mealType] = make([]uint32, 0)
				}
				menuConf[mealType] = append(menuConf[mealType], value.ID)
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

func GenerateMenuTypeListTableData(menuTypeList []*model.MenuType, dishTypeMap map[uint32]*model.DishType) []*dto.TableRowInfo {
	dataList := make([]*dto.TableRowInfo, 0)
	for _, menuType := range menuTypeList {
		row := &dto.TableRowInfo{DataMap: make(map[string]*dto.TableRowColumnInfo)}
		row.DataMap[IndexID] = &dto.TableRowColumnInfo{ID: menuType.ID, Value: fmt.Sprintf("%v", menuType.ID)}
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

	for mealType := enum.MealUnknown + 1; mealType < enum.MealALL; mealType++ {
		mealContent, ok := menuContentMap[mealType]
		if ok == false {
			continue
		}
		mealHead := &dto.TableColumnInfo{Name: enum.GetMealName(mealType), DataIndex: enum.GetMealKey(mealType),
			Hide: false, Children: make([]*dto.TableColumnInfo, 0)}
		for dishType := range mealContent {
			dishHead := &dto.TableColumnInfo{Name: dishTypeMap[dishType].DishTypeName,
				DataIndex: enum.GetMealKey(mealType) + fmt.Sprintf("_%v", dishTypeMap[dishType].ID), Hide: false}
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
	row.DataMap[IndexID] = &dto.TableRowColumnInfo{ID: menuType.ID, Value: fmt.Sprintf("%v", menuType.ID)}
	row.DataMap[IndexMenuTypeName] = &dto.TableRowColumnInfo{Value: menuType.MenuTypeName}
	for mealType, mealContent := range menuContentMap {
		for dishType, dishNumber := range mealContent {
			row.DataMap[enum.GetMealKey(mealType)+fmt.Sprintf("_%v", dishTypeMap[dishType].ID)] =
				&dto.TableRowColumnInfo{Value: fmt.Sprintf("%v", dishNumber)}
		}
	}
	dataList = append(dataList, row)
	return dataList
}

func ParseMenuTypeDetailData(row *dto.TableRowInfo, menuTypeID uint32, menuTypeName string) *model.MenuType {
	menuType := &model.MenuType{ID: menuTypeID, MenuTypeName: menuTypeName}
	menuConf := make(map[uint8]map[uint32]int32)
	for key, info := range row.DataMap {
		keys := strings.Split(key, IndexDelimiter)
		if len(keys) <= 1 {
			continue
		}
		mealType := enum.GetMealType(keys[0])
		if mealType == enum.MealUnknown {
			logger.Warn(TableGeneratorLogTag, "MenuTypeDetail Unknown Meal|Key:%v|Value:%#v", key, info)
			continue
		}
		if _, ok := menuConf[mealType]; ok == false {
			menuConf[mealType] = make(map[uint32]int32)
		}
		dishType, _ := strconv.ParseInt(keys[1], 10, 32)
		dishNumber, _ := strconv.ParseInt(info.Value, 10, 32)
		menuConf[mealType][uint32(dishType)] = int32(dishNumber)
	}
	menuType.FromMenuConfig(menuConf)
	return menuType
}

func GenerateWeekMenuListTableHead() []*dto.TableColumnInfo {
	head := make([]*dto.TableColumnInfo, 0)
	head = append(head, &dto.TableColumnInfo{Name: "菜单ID", DataIndex: IndexID, Hide: false})
	head = append(head, &dto.TableColumnInfo{Name: "菜单日期", DataIndex: IndexMenuDate, Hide: false})
	for _, dayName := range weekDays {
		for i := enum.MealUnknown + 1; i < enum.MealALL; i++ {
			head = append(head, &dto.TableColumnInfo{Name: enum.GetMealName(i) + dayName,
				DataIndex: GenerateDayMealIndex(i, dayName), Hide: false})
		}
	}

	return head
}

func GenerateWeekMenuListTableData(menuList []*model.WeekMenu, dishIDMap map[uint32]*model.Dish) []*dto.TableRowInfo {
	dataList := make([]*dto.TableRowInfo, 0)
	for _, menu := range menuList {
		row := &dto.TableRowInfo{DataMap: make(map[string]*dto.TableRowColumnInfo)}
		row.DataMap[IndexID] = &dto.TableRowColumnInfo{ID: menu.ID, Value: fmt.Sprintf("%v", menu.ID)}
		startDate := menu.MenuStartDate.Format("01-02")
		endDate := menu.MenuStartDate.Add(time.Hour * 24 * 7).Format("01-02")
		row.DataMap[IndexMenuDate] = &dto.TableRowColumnInfo{Value: fmt.Sprintf("%v~%v", startDate, endDate)}
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
				row.DataMap[GenerateDayMealIndex(mealType, weekDays[dayIndex%len(weekDays)])] =
					&dto.TableRowColumnInfo{Value: mealStr[:len(mealStr)-1]}
			}
		}
	}
	return dataList
}

func GenerateWeekMenuDetailTableHead(menu *model.WeekMenu, dishIDMap map[uint32]*model.Dish,
	dishTypeMap map[uint32]*model.DishType) ([]*dto.TableColumnInfo, []*dto.TableRowInfo) {
	head := make([]*dto.TableColumnInfo, 0)
	dataList := make([]*dto.TableRowInfo, 0)
	head = append(head, &dto.TableColumnInfo{Name: "日期", DataIndex: IndexMenuDate, Hide: false})
	for _, dayName := range weekDays {
		for i := enum.MealUnknown + 1; i < enum.MealALL; i++ {
			head = append(head, &dto.TableColumnInfo{Name: enum.GetMealName(i) + dayName,
				DataIndex: GenerateDayMealIndex(i, dayName), Hide: false})
		}
	}

	return head, dataList
}
