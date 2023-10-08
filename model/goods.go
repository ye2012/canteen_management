package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		"batch_size", "batch_unit", "quantity"}
)

type Goods struct {
	ID           uint32    `json:"id"`
	Name         string    `json:"name"`
	GoodsTypeID  uint32    `json:"goods_type_id"`
	StoreTypeID  uint32    `json:"store_type_id"`
	Picture      string    `json:"picture"`
	BatchSize    float64   `json:"batch_size"`
	BatchUnit    string    `json:"batch_unit"`
	Price        float64   `json:"price"`
	PriceContent string    `json:"price_content"`
	Quantity     float64   `json:"quantity"`
	CreateAt     time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}

func (g *Goods) FromGoodsPrice(priceConf map[uint8]float64) error {
	contentMap := make(map[uint8]uint64)
	for id, price := range priceConf {
		contentMap[id] = uint64(price * 100)
	}
	contentStr, err := json.Marshal(contentMap)
	if err != nil {
		logger.Warn(goodsLogTag, "FromGoodsPrice Failed|Err:%v", err)
		return err
	}
	g.PriceContent = string(contentStr)
	return nil
}

func (g *Goods) ToGoodsPrice() map[uint8]float64 {
	contentMap := make(map[uint8]uint64)
	err := json.Unmarshal([]byte(g.PriceContent), &contentMap)
	if err != nil {
		logger.Warn(goodsLogTag, "ToGoodsPrice Failed|Err:%v", err)
		return nil
	}
	retMap := make(map[uint8]float64)
	for id, price := range contentMap {
		retMap[id] = float64(price) / 100
	}
	return retMap
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

func (gm *GoodsModel) GetAllGoods() ([]*Goods, error) {
	retList, err := utils.SqlQuery(gm.sqlCli, goodsTable, &Goods{}, "")
	if err != nil {
		logger.Warn(goodsLogTag, "GetGoods Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*Goods), nil
}

func (gm *GoodsModel) GetGoods(goodsType, storeType uint32, page, pageSize int32) ([]*Goods, error) {
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
	condition += " ORDER BY `id` DESC LIMIT ?,? "
	params = append(params, (page-1)*pageSize, pageSize)
	retList, err := utils.SqlQuery(gm.sqlCli, goodsTable, &Goods{}, condition, params...)
	if err != nil {
		logger.Warn(goodsLogTag, "GetGoods Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*Goods), nil
}

func (gm *GoodsModel) GetGoodsCount(goodsType, storeType uint32) (int32, error) {
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
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM %v %v", goodsTable, condition)
	row := gm.sqlCli.QueryRow(sqlStr, params...)
	var count int32
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (gm *GoodsModel) GetGoodsByID(id uint32) (*Goods, error) {
	goods := &Goods{}
	err := utils.SqlQueryRow(gm.sqlCli, goodsTable, goods, " WHERE `id` = ? ", id)
	if err != nil {
		logger.Warn(goodsLogTag, "GetGoods Failed|Err:%v", err)
		return nil, err
	}

	return goods, nil
}
func (gm *GoodsModel) BatchAddQuantity(updateList []*Goods) (err error) {
	return gm.BatchAddQuantityWithTx(nil, updateList)
}

func (gm *GoodsModel) BatchAddQuantityWithTx(tx *sql.Tx, updateList []*Goods) (err error) {
	daoList := make([]interface{}, 0)
	for _, updateInfo := range updateList {
		daoList = append(daoList, updateInfo)
	}
	if tx != nil {
		err = utils.SqlBatchAdd(tx, goodsTable, daoList, "id", "quantity")
	} else {
		err = utils.SqlBatchAdd(gm.sqlCli, goodsTable, daoList, "id", "quantity")
	}
	if err != nil {
		logger.Warn(goodsLogTag, "BatchUpdateQuantity Failed|Err:%v", err)
		return err
	}
	return nil
}

func (gm *GoodsModel) UpdateGoodsInfo(dao *Goods) error {
	err := utils.SqlUpdateWithUpdateTags(gm.sqlCli, goodsTable, dao, "id", goodsUpdateTags...)
	if err != nil {
		logger.Warn(goodsLogTag, "UpdateGoodsInfo Failed|Err:%v", err)
		return err
	}
	return nil
}

func (gm *GoodsModel) UpdateGoodsPriceInfo(id uint32, price float64, priceMap map[uint8]float64) error {
	dao := &Goods{ID: id, Price: price}
	err := dao.FromGoodsPrice(priceMap)
	if err != nil {
		logger.Warn(goodsLogTag, "UpdatePriceInfo FromGoodsPrice  Failed|Err:%v", err)
		return err
	}
	err = utils.SqlUpdateWithUpdateTags(gm.sqlCli, goodsTable, dao, "id", "price", "price_content")
	if err != nil {
		logger.Warn(goodsLogTag, "UpdateGoodsPriceInfo Failed|Err:%v", err)
		return err
	}
	return nil
}
