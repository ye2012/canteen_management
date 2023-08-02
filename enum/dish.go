package enum

import "fmt"

type MealType = uint8

const (
	MealUnknown MealType = iota
	MealBreakfast
	MealLunch
	MealDinner
	MealALL
)

const (
	MealBreakfastName = "早餐"
	MealLunchName     = "午餐"
	MealDinnerName    = "晚餐"

	MealBreakfastKey = "Breakfast"
	MealLunchKey     = "Lunch"
	MealDinnerKey    = "Dinner"
)

var mealNameMap = map[MealType]string{
	MealBreakfast: MealBreakfastName,
	MealLunch:     MealLunchName,
	MealDinner:    MealDinnerName,
}

var mealKeyMap = map[MealType]string{
	MealBreakfast: MealBreakfastKey,
	MealLunch:     MealLunchKey,
	MealDinner:    MealDinnerKey,
}

var mealTypeMap = map[string]MealType{
	MealBreakfastName: MealBreakfast,
	MealLunchName:     MealLunch,
	MealDinnerName:    MealDinner,
}

func GetMealName(timeType MealType) string {
	name, ok := mealNameMap[timeType]
	if ok {
		return name
	}
	return fmt.Sprintf("Unknown Meal:%v", timeType)
}

func GetMealKey(timeType MealType) string {
	name, ok := mealKeyMap[timeType]
	if ok {
		return name
	}
	return fmt.Sprintf("Unknown Meal:%v", timeType)
}

func GetMealType(name string) MealType {
	timeType, ok := mealTypeMap[name]
	if ok {
		return timeType
	}
	return MealUnknown
}
