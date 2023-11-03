package service

import (
	"database/sql"
	"fmt"
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
}

func NewDishService(sqlCli *sql.DB) *DishService {
	dishModel := model.NewDishesModelWithDB(sqlCli)
	dishTypeModel := model.NewDishTypeModelWithDB(sqlCli)
	return &DishService{
		dishModel:     dishModel,
		dishTypeModel: dishTypeModel,
	}
}

func (ds *DishService) Init() error {
	return nil
}

func (ds *DishService) GetDishIDMap() (map[uint32]*model.Dish, error) {
	dishList, err := ds.dishModel.GetDishes(0, 0, 100000)
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishIDMap GetDishes Failed|Err:%v", err)
		return nil, err
	}
	idMap := make(map[uint32]*model.Dish)
	for _, dish := range dishList {
		idMap[dish.ID] = dish
	}
	return idMap, nil
}

func (ds *DishService) GetDishMap() map[uint32][]*model.Dish {
	dishList, err := ds.dishModel.GetDishes(0, 0, 100000)
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishMap GetDishes Failed|Err:%v", err)
		return nil
	}

	typeMap := make(map[uint32][]*model.Dish)
	for _, dish := range dishList {
		_, ok := typeMap[dish.DishType]
		if ok == false {
			typeMap[dish.DishType] = make([]*model.Dish, 0)
		}
		typeMap[dish.DishType] = append(typeMap[dish.DishType], dish)
	}
	return typeMap
}

func (ds *DishService) GetDishList(dishType uint32, page, pageSize int32) ([]*model.Dish, int32, error) {
	dishList, err := ds.dishModel.GetDishes(dishType, page, pageSize)
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishList Failed|Err:%v", err)
		return nil, 0, err
	}

	dishCount, err := ds.dishModel.GetDishesCount(dishType)
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishesCount Failed|Err:%v", err)
		return nil, 0, err
	}
	return dishList, dishCount, nil
}

func (ds *DishService) AddDish(dish *model.Dish) error {
	dishDao, err := ds.dishModel.GetDishByName(dish.DishName)
	if err != nil {
		logger.Warn(dishServiceLogTag, "AddDish GetDishByName Failed|Err:%v", err)
		return err

	}
	if dishDao != nil {
		return fmt.Errorf("已有同名菜品，无法添加")
	}

	err = ds.dishModel.Insert(dish)
	if err != nil {
		logger.Warn(dishServiceLogTag, "Insert Dish Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ds *DishService) ModifyDish(dish *model.Dish) error {
	err := ds.dishModel.UpdateDish(dish)
	if err != nil {
		logger.Warn(dishServiceLogTag, "Update Dish Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ds *DishService) GetDishTypeMap() (map[uint32]*model.DishType, error) {
	dishTypeList, err := ds.dishTypeModel.GetDishTypes()
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishTypeMap Failed|Err:%v", err)
		return nil, err
	}

	retMap := make(map[uint32]*model.DishType)
	for _, typeInfo := range dishTypeList {
		retMap[typeInfo.ID] = typeInfo
	}
	return retMap, nil
}

func (ds *DishService) GetDishTypeList(masterType uint32, includeMaster bool, page, pageSize int32) ([]*model.DishType, int32, error) {
	if masterType == 0 && includeMaster {
		dishTypeList, err := ds.dishTypeModel.GetMasterDishTypes()
		if err != nil {
			logger.Warn(dishServiceLogTag, "GetMasterDishTypes Failed|Err:%v", err)
			return nil, 0, err
		}
		return dishTypeList, int32(len(dishTypeList)), err
	}

	dishTypeList, err := ds.dishTypeModel.GetDishTypesByMasterType(masterType, page, pageSize)
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishTypesByMasterType Failed|Err:%v", err)
		return nil, 0, err
	}

	dishTypeCount, err := ds.dishTypeModel.GetDishTypesCountByMasterType(masterType)
	if err != nil {
		logger.Warn(dishServiceLogTag, "GetDishTypesCountByMasterType Failed|Err:%v", err)
		return nil, 0, err
	}

	return dishTypeList, dishTypeCount, err
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

func (ds *DishService) DeleteDishType(dishTypeID uint32) error {
	count, err := ds.dishModel.GetDishesCount(dishTypeID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "DelDishType GetDishesCount Failed|Err:%v", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("该菜品类型下还有菜品，无法删除")
	}
	err = ds.dishTypeModel.DeleteDishType(dishTypeID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsTypesByID Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ds *DishService) RandDishByType(typeID uint32, number int) []*model.Dish {
	dishTypeMap := ds.GetDishMap()
	if dishTypeMap == nil {
		return nil
	}
	dishList, ok := dishTypeMap[typeID]
	if ok == false {
		return nil
	}

	dishLen := len(dishList)
	times := number/dishLen + 1

	retList := make([]*model.Dish, 0)
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

func (ds *DishService) RandDishIDByType(typeID uint32, number int) []uint32 {
	dishTypeMap := ds.GetDishMap()
	if dishTypeMap == nil {
		return nil
	}
	dishList, ok := dishTypeMap[typeID]
	if ok == false {
		return nil
	}

	dishLen := len(dishList)
	times := number/dishLen + 1

	retList := make([]uint32, 0)
	for curRound := 0; curRound < times && number > 0; curRound++ {
		randList := rand.Perm(dishLen)
		for _, randIndex := range randList {
			retList = append(retList, dishList[randIndex%dishLen].ID)
			number--
			if number <= 0 {
				break
			}
		}
	}

	return retList
}
