package service

import (
	"database/sql"
	"math/rand"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	dishServiceLogTag = "DishService"
)

type DishService struct {
	dishModel     *model.DishesModel
	dishTypeModel *model.DishTypeModel

	dishTypeMap map[uint32][]*model.Dish // type => dish list
	dishIDMap   map[uint32]*model.Dish
}

func NewDishService(sqlCli *sql.DB) *DishService {
	dishModel := model.NewDishesModelWithDB(sqlCli)
	dishTypeModel := model.NewDishTypeModelWithDB(sqlCli)
	return &DishService{
		dishModel:     dishModel,
		dishTypeModel: dishTypeModel,
		dishTypeMap:   make(map[uint32][]*model.Dish),
		dishIDMap:     make(map[uint32]*model.Dish),
	}
}

func (ds *DishService) Init() error {
	dishList, err := ds.dishModel.GetDishes(0)
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishList Failed|Err:%v", err)
		return err
	}

	for _, dish := range dishList {
		ds.addDishCache(dish)
	}
	return nil
}

func (ds *DishService) addDishCache(dish *model.Dish) {
	_, ok := ds.dishTypeMap[dish.DishType]
	if ok == false {
		ds.dishTypeMap[dish.DishType] = make([]*model.Dish, 0)
	}
	ds.dishTypeMap[dish.DishType] = append(ds.dishTypeMap[dish.DishType], dish)

	_, ok = ds.dishIDMap[dish.ID]
	if ok {
		logger.Warn(dishServiceLogTag, "DishID Complex|ID:%v", dish.ID)
	}
	ds.dishIDMap[dish.ID] = dish
}

func (ds *DishService) GetDishIDMap() map[uint32]*model.Dish {
	return ds.dishIDMap
}

func (ds *DishService) updateDishCache(dish *model.Dish) {
	_, ok := ds.dishIDMap[dish.ID]
	if ok == false {
		logger.Warn(dishServiceLogTag, "DishID Not Found|Type:%v", dish.DishType)
		return
	}
	preType := ds.dishIDMap[dish.ID].DishType
	ds.dishIDMap[dish.ID] = dish

	if preType != dish.DishType {
		_, ok = ds.dishTypeMap[preType]
		if ok {
			for index, dishDao := range ds.dishTypeMap[preType] {
				if dishDao.ID == dish.ID {
					ds.dishTypeMap[preType] = append(ds.dishTypeMap[preType][:index], ds.dishTypeMap[preType][index+1:]...)
				}
			}
		}
	}

	_, ok = ds.dishTypeMap[dish.DishType]
	if ok == false {
		logger.Warn(dishServiceLogTag, "DishType Not Found|Type:%v", dish.DishType)
		return
	}
	for index, dishDao := range ds.dishTypeMap[dish.DishType] {
		if dish.ID == dishDao.ID {
			ds.dishTypeMap[dish.DishType][index] = dish
			return
		}
	}

}

func (ds *DishService) GetDishList(dishType uint32) ([]*model.Dish, error) {
	dishList, err := ds.dishModel.GetDishes(dishType)
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishList Failed|Err:%v", err)
		return nil, err
	}
	return dishList, nil
}

func (ds *DishService) AddDish(dish *model.Dish) error {
	err := ds.dishModel.Insert(dish)
	if err != nil {
		logger.Warn(dishServiceLogTag, "Insert Dish Failed|Err:%v", err)
		return err
	}
	ds.addDishCache(dish)
	return nil
}

func (ds *DishService) ModifyDish(dish *model.Dish) error {
	err := ds.dishModel.UpdateDish(dish)
	if err != nil {
		logger.Warn(dishServiceLogTag, "Update Dish Failed|Err:%v", err)
		return err
	}
	ds.updateDishCache(dish)
	return nil
}

func (ds *DishService) GetDishTypeList() ([]*model.DishType, error) {
	dishTypeList, err := ds.dishTypeModel.GetDishTypes()
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishTypes Failed|Err:%v", err)
		return nil, err
	}
	return dishTypeList, err
}

func (ds *DishService) AddDishType(dishType *model.DishType) error {
	err := ds.dishTypeModel.Insert(dishType)
	if err != nil {
		logger.Warn(dishServiceLogTag, "Insert DishType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ds *DishService) ModifyDishType(dishType *model.DishType) error {
	err := ds.dishTypeModel.UpdateDishType(dishType)
	if err != nil {
		logger.Warn(dishServiceLogTag, "Update DishType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ds *DishService) RandDishByType(typeID uint32, number int) []*model.Dish {
	retList := make([]*model.Dish, 0)
	dishList, ok := ds.dishTypeMap[typeID]
	if ok == false {
		return retList
	}

	dishLen := len(dishList)
	times := number/dishLen + 1

	for curRound := 0; curRound < times && number > 0; curRound++ {
		randList := rand.Perm(dishLen)
		for _, randIndex := range randList {
			retList = append(retList, dishList[randIndex%dishLen])
			number--
			if number <= 0 {
				break
			}
		}
	}

	return retList
}
