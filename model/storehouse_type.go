package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	storehouseTypeTable = "storehouse_type"

	storehouseTypeLogTag = "StorehouseModel"
)

var (
	storeTypeUpdateTags = []string{"store_type_name"}
)

type StorehouseType struct {
	ID            uint32    `json:"id"`
	StoreTypeName string    `json:"store_type_name"`
	CreateAt      time.Time `json:"created_at"`
	UpdateAt      time.Time `json:"updated_at"`
}

type StorehouseTypeModel struct {
	sqlCli *sql.DB
}

func NewStorehouseTypeModelWithDB(sqlCli *sql.DB) *StorehouseTypeModel {
	return &StorehouseTypeModel{
		sqlCli: sqlCli,
	}
}

func (stm *StorehouseTypeModel) Insert(dao *StorehouseType) error {
	id, err := utils.SqlInsert(stm.sqlCli, storehouseTypeTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(storehouseTypeLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (stm *StorehouseTypeModel) GetStorehouseTypes() ([]*StorehouseType, error) {
	retList, err := utils.SqlQuery(stm.sqlCli, storehouseTypeTable, &StorehouseType{}, "")
	if err != nil {
		logger.Warn(storehouseTypeLogTag, "GetStorehouseTypes Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*StorehouseType), nil
}

func (stm *StorehouseTypeModel) UpdateStorehouseType(dao *StorehouseType) error {
	err := utils.SqlUpdateWithUpdateTags(stm.sqlCli, storehouseTypeTable, dao, "id", storeTypeUpdateTags...)
	if err != nil {
		logger.Warn(storehouseTypeLogTag, "UpdateStorehouseType Failed|Err:%v", err)
		return err
	}
	return nil
}
