package enum

import "fmt"

type MealType = uint8

const (
	MealUnknown MealType = iota
	MealBreakfast
	MealLunch
	MealDinner
)

const (
	MealBreakfastName = "Breakfast"
	MealLunchName     = "Lunch"
	MealDinnerName    = "Dinner"
)

var mealNameMap = map[MealType]string{
	MealBreakfast: MealBreakfastName,
	MealLunch:     MealLunchName,
	MealDinner:    MealDinnerName,
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

func GetMealType(name string) MealType {
	timeType, ok := mealTypeMap[name]
	if ok {
		return timeType
	}
	return MealUnknown
}
