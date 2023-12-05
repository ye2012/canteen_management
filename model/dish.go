package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	dishLogTag = "DishModel"

	dishTable = "dishes"
)

var (
	dishUpdateTags = []string{"picture", "dish_name", "dish_type", "price", "material"}
)

type Dish struct {
	ID       uint32    `json:"id"`
	DishName string    `json:"dish_name"`
	DishType uint32    `json:"dish_type"`
	Picture  string    `json:"picture"`
	Price    float64   `json:"price"`
	Material string    `json:"material"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type DishesModel struct {
	sqlCli *sql.DB
}

func NewDishesModelWithDB(sqlCli *sql.DB) *DishesModel {
	return &DishesModel{
		sqlCli: sqlCli,
	}
}

func NewDishesModel(config utils.Config) (*DishesModel, error) {
	dbClient, err := utils.NewMysqlClient(config)
	if err != nil {
		logger.Error(dishLogTag, "New Service Fail|Err:%v", err)
		return nil, err
	}

	return &DishesModel{
		sqlCli: dbClient,
	}, nil
}

func (dm *DishesModel) Insert(dao *Dish) error {
	id, err := utils.SqlInsert(dm.sqlCli, dishTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(dishLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (dm *DishesModel) BatchInsert(daoList []*Dish) error {
	err := utils.SqlInsertBatch(dm.sqlCli, dishTable, daoList, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(dishLogTag, "BatchInsert Failed|Err:%v", err)
		return err
	}
	return nil
}

func (dm *DishesModel) GetDishByCondition(condition string, params ...interface{}) ([]*Dish, error) {
	retList, err := utils.SqlQuery(dm.sqlCli, dishTable, &Dish{}, condition, params...)
	if err != nil {
		logger.Warn(dishLogTag, "GetDishes Failed|Err:%v", err)
		return nil, err
	}
	return retList.([]*Dish), nil
}

func (dm *DishesModel) GetDishCountByCondition(condition string, params ...interface{}) (int32, error) {
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM %v %v", dishTable, condition)
	row := dm.sqlCli.QueryRow(sqlStr, params...)
	var count int32
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (dm *DishesModel) GetDishByName(dishName string) (*Dish, error) {
	condition := " WHERE `dish_name` = ? "
	dishList, err := dm.GetDishByCondition(condition, dishName)
	if err != nil {
		logger.Warn(dishLogTag, "GetDishByName Failed|Err:%v", err)
		return nil, err
	}
	if len(dishList) > 0 {
		return dishList[0], nil
	}
	return nil, nil
}

func (dm *DishesModel) GetDishByNameList(dishNameList []string) ([]*Dish, error) {
	if len(dishNameList) == 0 {
		return make([]*Dish, 0), nil
	}
	params := make([]interface{}, 0, len(dishNameList))
	for _, dishName := range dishNameList {
		params = append(params, dishName)
	}
	condition := fmt.Sprintf(" WHERE `dish_name` in (%v) ", strings.Repeat(",?", len(dishNameList))[1:])
	dishList, err := dm.GetDishByCondition(condition, params...)
	if err != nil {
		logger.Warn(dishLogTag, "GetDishByName Failed|Err:%v", err)
		return nil, err
	}
	return dishList, nil
}

func (dm *DishesModel) GetDishes(dishType uint32, page, pageSize int32) ([]*Dish, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if dishType > 0 {
		condition += " AND `dish_type` = ? "
		params = append(params, dishType)
	}

	if page >= 1 {
		condition += " ORDER BY `id` ASC LIMIT ?,? "
		params = append(params, (page-1)*pageSize, pageSize)
	}

	return dm.GetDishByCondition(condition, params...)
}

func (dm *DishesModel) GetDishesCount(dishType uint32) (int32, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if dishType > 0 {
		condition += " AND `dish_type` = ? "
		params = append(params, dishType)
	}
	return dm.GetDishCountByCondition(condition, params...)
}

func (dm *DishesModel) UpdateDish(dao *Dish) error {
	err := utils.SqlUpdateWithUpdateTags(dm.sqlCli, dishTable, dao, "id", dishUpdateTags...)
	if err != nil {
		logger.Warn(dishLogTag, "UpdateDish Failed|Err:%v", err)
		return err
	}
	return nil
}

func (dm *DishesModel) DeleteDish(dishID uint32) error {
	sqlStr := fmt.Sprintf(" DELETE FROM %v WHERE `id` = ? ", dishTable)
	_, err := dm.sqlCli.Exec(sqlStr, dishID)
	if err != nil {
		logger.Warn(dishLogTag, "DeleteDish Failed|Err:%v", err)
		return err
	}

	return nil
}
