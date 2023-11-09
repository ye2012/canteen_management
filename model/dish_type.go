package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	dishTypeLogTag = "DishType"

	dishTypeTable = "dish_type"
)

var dishTypeUpdateTags = []string{"dish_type_name", "master_type"}

type DishType struct {
	ID           uint32    `json:"id"`
	DishTypeName string    `json:"dish_type_name"`
	MasterType   uint32    `json:"master_type"`
	CreateAt     time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}

type DishTypeList []*DishType

func (dtl DishTypeList) Len() int {
	return dtl.Len()
}

func (dtl DishTypeList) Less(i, j int) bool {
	return dtl[i].ID > dtl[j].ID
}

func (dtl DishTypeList) Swap(i, j int) {
	dtl[i], dtl[j] = dtl[j], dtl[i]
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

func (dtm *DishTypeModel) GetMasterDishTypes() ([]*DishType, error) {
	condition := " WHERE `master_type` = 0 "
	retList, err := dtm.GetDishTypesByCondition(condition)
	if err != nil {
		logger.Warn(dishTypeLogTag, "GetMasterDishTypes Failed|Err:%v", err)
		return nil, err
	}
	return retList, nil
}

func (dtm *DishTypeModel) GetDishTypesByMasterType(masterTypeID uint32, page, pageSize int32) ([]*DishType, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if masterTypeID > 0 {
		condition += " AND `master_type` = ? "
		params = append(params, masterTypeID)
	} else {
		condition += " AND `master_type` > 0 "
	}
	if page >= 1 {
		condition += " ORDER BY `id` ASC LIMIT ?,? "
		params = append(params, (page-1)*pageSize, pageSize)
	}

	retList, err := dtm.GetDishTypesByCondition(condition, params...)
	if err != nil {
		logger.Warn(dishTypeLogTag, "GetDishTypesByMasterType Failed|Err:%v", err)
		return nil, err
	}
	return retList, nil
}

func (dtm *DishTypeModel) GetDishTypesCountByMasterType(masterTypeID uint32) (int32, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if masterTypeID > 0 {
		condition += " AND `master_type` = ? "
		params = append(params, masterTypeID)
	} else {
		condition += " AND `master_type` > 0 "
	}

	count, err := dtm.GetDishTypesCountByCondition(condition, params...)
	if err != nil {
		logger.Warn(dishTypeLogTag, "GetDishTypesCountByCondition Failed|Err:%v", err)
		return 0, err
	}
	return count, nil
}

func (dtm *DishTypeModel) GetDishTypesByCondition(condition string, params ...interface{}) ([]*DishType, error) {
	retList, err := utils.SqlQuery(dtm.sqlCli, dishTypeTable, &DishType{}, condition, params...)
	if err != nil {
		logger.Warn(dishTypeLogTag, "GetDishTypes Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*DishType), nil
}

func (dtm *DishTypeModel) GetDishTypesCountByCondition(condition string, params ...interface{}) (int32, error) {
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM %v %v", dishTypeTable, condition)
	row := dtm.sqlCli.QueryRow(sqlStr, params...)
	var count int32
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (dtm *DishTypeModel) UpdateDishType(dao *DishType) error {
	err := utils.SqlUpdateWithUpdateTags(dtm.sqlCli, dishTypeTable, dao, "id", dishTypeUpdateTags...)
	if err != nil {
		logger.Warn(dishTypeLogTag, "UpdateDishType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (dtm *DishTypeModel) DeleteDishType(id uint32) error {
	sqlStr := fmt.Sprintf(" DELETE FROM %v WHERE `id` = ? ", dishTypeTable)
	_, err := dtm.sqlCli.Exec(sqlStr, id)
	if err != nil {
		logger.Warn(dishTypeLogTag, "DeleteDishType Failed|Err:%v", err)
		return err
	}

	return nil
}
