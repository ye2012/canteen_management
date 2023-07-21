package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	dishTypeLogTag = "DishType"

	dishTypeTable = "dish_type"
)

var dishTypeUpdateTags = []string{"dish_type_name"}

type DishType struct {
	ID           uint32    `json:"id"`
	DishTypeName string    `json:"dish_type_name"`
	CreateAt     time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}

type DishTypeModel struct {
	sqlCli *sql.DB
}

func NewDishTypeModelWithDB(sqlCli *sql.DB) *DishTypeModel {
	return &DishTypeModel{
		sqlCli: sqlCli,
	}
}

func (dtm *DishTypeModel) Insert(dao *DishType) error {
	id, err := utils.SqlInsert(dtm.sqlCli, dishTypeTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(dishTypeLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (dtm *DishTypeModel) GetDishTypes() ([]*DishType, error) {
	retList, err := utils.SqlQuery(dtm.sqlCli, dishTypeTable, &DishType{}, "")
	if err != nil {
		logger.Warn(dishTypeLogTag, "GetDishTypes Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*DishType), nil
}

func (dtm *DishTypeModel) UpdateDishType(dao *DishType) error {
	err := utils.SqlUpdateWithUpdateTags(dtm.sqlCli, dishTypeTable, dao, "id", dishTypeUpdateTags...)
	if err != nil {
		logger.Warn(dishTypeLogTag, "UpdateDishType Failed|Err:%v", err)
		return err
	}
	return nil
}
