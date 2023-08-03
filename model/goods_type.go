package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	goodsTypeTable = "goods_type"

	goodsTypeLogTag = "GoodsType"
)

var (
	goodsTypeUpdateTags = []string{"goods_type_name", "discount"}
)

type GoodsType struct {
	ID            uint32    `json:"id"`
	GoodsTypeName string    `json:"goods_type_name"`
	Discount      float64   `json:"discount"`
	CreateAt      time.Time `json:"created_at"`
	UpdateAt      time.Time `json:"updated_at"`
}

type GoodsTypeModel struct {
	sqlCli *sql.DB
}

func NewGoodsTypeModelWithDB(sqlCli *sql.DB) *GoodsTypeModel {
	return &GoodsTypeModel{
		sqlCli: sqlCli,
	}
}

func (gtm *GoodsTypeModel) Insert(dao *GoodsType) error {
	id, err := utils.SqlInsert(gtm.sqlCli, goodsTypeTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(goodsTypeLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (gtm *GoodsTypeModel) GetGoodsTypes() ([]*GoodsType, error) {
	retList, err := utils.SqlQuery(gtm.sqlCli, goodsTypeTable, &GoodsType{}, "")
	if err != nil {
		logger.Warn(goodsTypeLogTag, "GetGoodsTypes Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*GoodsType), nil
}

func (gtm *GoodsTypeModel) UpdateGoodsType(dao *GoodsType) error {
	err := utils.SqlUpdateWithUpdateTags(gtm.sqlCli, goodsTypeTable, dao, "id", goodsTypeUpdateTags...)
	if err != nil {
		logger.Warn(goodsTypeLogTag, "UpdateGoodsType Failed|Err:%v", err)
		return err
	}
	return nil
}
