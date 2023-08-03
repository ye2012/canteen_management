package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	goodsTable = "goods"

	goodsLogTag = "Goods"
)

var (
	goodsUpdateTags = []string{"name", "goods_type_id", "store_type_id", "picture",
		"batch_size", "batch_unit", "price", "quantity"}
)

type Goods struct {
	ID          uint32    `json:"id"`
	Name        string    `json:"name"`
	GoodsTypeID uint32    `json:"goods_type_id"`
	StoreTypeID uint32    `json:"store_type_id"`
	Picture     string    `json:"picture"`
	BatchSize   float64   `json:"batch_size"`
	BatchUnit   string    `json:"batch_unit"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

type GoodsModel struct {
	sqlCli *sql.DB
}

func NewGoodsModelWithDB(sqlCli *sql.DB) *GoodsModel {
	return &GoodsModel{
		sqlCli: sqlCli,
	}
}

func (gm *GoodsModel) Insert(dao *Goods) error {
	id, err := utils.SqlInsert(gm.sqlCli, goodsTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(goodsLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (gm *GoodsModel) GetGoods(goodsType, storeType uint32) ([]*Goods, error) {
	condition := " WHERE 1=1 "
	var params []interface{}
	if goodsType > 0 {
		condition += " AND `goods_type_id` = ? "
		params = append(params, goodsType)
	}
	if storeType > 0 {
		condition += " AND `store_type_id` = ? "
		params = append(params, storeType)
	}
	retList, err := utils.SqlQuery(gm.sqlCli, goodsTable, &Goods{}, condition, params...)
	if err != nil {
		logger.Warn(goodsLogTag, "GetGoods Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*Goods), nil
}

func (gm *GoodsModel) UpdateGoodsInfo(dao *Goods) error {
	err := utils.SqlUpdateWithUpdateTags(gm.sqlCli, goodsTable, dao, "id", goodsUpdateTags...)
	if err != nil {
		logger.Warn(goodsLogTag, "UpdateGoodsInfo Failed|Err:%v", err)
		return err
	}
	return nil
}
