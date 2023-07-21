package model

import (
	"database/sql"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
	"time"
)

const (
	storehouseTable = "StorehouseDetail"

	storehouseLogTag = "StorehouseModel"
)

var (
	storeDetailUpdateTags = []string{"quantity"}
)

type StorehouseDetail struct {
	ID          uint32    `json:"id"`
	StoreTypeID uint32    `json:"store_type_id"`
	GoodsID     uint32    `json:"goods_id"`
	Quantity    float64   `json:"quantity"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

type StorehouseModel struct {
	sqlCli *sql.DB
}

func NewStorehouseModelWithDB(sqlCli *sql.DB) *StorehouseModel {
	return &StorehouseModel{
		sqlCli: sqlCli,
	}
}

func (sm *StorehouseModel) Insert(dao *StorehouseDetail) error {
	id, err := utils.SqlInsert(sm.sqlCli, storehouseTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(storehouseLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (sm *StorehouseModel) GetStorehouseGoods(storeType uint32) ([]*StorehouseDetail, error) {
	condition, params := " WHERE 1=1 ", make([]interface{}, 0)
	if storeType > 0 {
		condition += " AND `store_type_id` = ? "
		params = append(params, storeType)
	}
	retList, err := utils.SqlQuery(sm.sqlCli, storehouseTable, &StorehouseDetail{}, condition, params...)
	if err != nil {
		logger.Warn(storehouseLogTag, "GetStorehouseTypes Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*StorehouseDetail), nil
}

func (sm *StorehouseModel) UpdateStorehouse(dao *StorehouseDetail) error {
	err := utils.SqlUpdateWithUpdateTags(sm.sqlCli, storehouseTable, dao, "id", storeDetailUpdateTags...)
	if err != nil {
		logger.Warn(storehouseLogTag, "UpdateStorehouse Failed|Err:%v", err)
		return err
	}
	return nil
}
