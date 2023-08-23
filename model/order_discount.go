package model

import (
	"database/sql"
	"encoding/json"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
	"time"
)

const (
	orderDiscountTable = "order_discount"

	orderDiscountLogTag = "OrderDiscount"
)

var (
	orderDiscountUpdateTags = []string{"discount_type_name", "discount_conf"}
)

type OrderDiscount struct {
	ID               uint32            `json:"id"`
	DiscountTypeName string            `json:"discount_type_name"`
	DiscountConf     string            `json:"discount_conf"`
	CreateAt         time.Time         `json:"created_at"`
	UpdateAt         time.Time         `json:"updated_at"`
	discountMap      map[uint8]float64 `json:"-"`
}

func (or *OrderDiscount) ConvertDiscount() error {
	or.discountMap = make(map[uint8]float64)
	if or.DiscountConf != "" {
		err := json.Unmarshal([]byte(or.DiscountConf), &or.discountMap)
		if err != nil {
			logger.Warn(orderDiscountLogTag, "ConvertDiscount Failed|Err:%v", err)
			return err
		}
	}
	return nil
}

func (or *OrderDiscount) GetMealDiscount(mealType uint8) float64 {
	if or.discountMap == nil {
		or.ConvertDiscount()
	}
	discount, ok := or.discountMap[mealType]
	if ok {
		return discount
	}
	return 0
}

type OrderDiscountModel struct {
	sqlCli *sql.DB
}

func NewOrderDiscountModel(sqlCli *sql.DB) *OrderDiscountModel {
	return &OrderDiscountModel{
		sqlCli: sqlCli,
	}
}

func (odm *OrderDiscountModel) Insert(dao *OrderDiscount) error {
	id, err := utils.SqlInsert(odm.sqlCli, orderDiscountTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(orderDiscountLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (odm *OrderDiscountModel) UpdateDiscountType(dao *OrderDiscount) error {
	err := utils.SqlUpdateWithUpdateTags(odm.sqlCli, orderDiscountTable, dao, "id", orderDiscountUpdateTags...)
	if err != nil {
		logger.Warn(orderDiscountLogTag, "UpdateDiscount Failed|Err:%v", err)
		return err
	}
	return nil
}

func (odm *OrderDiscountModel) GetDiscountByID(id uint8) (*OrderDiscount, error) {
	retInfo := &OrderDiscount{}
	err := utils.SqlQueryRow(odm.sqlCli, orderDiscountTable, retInfo, " WHERE `id`=? ", id)
	if err != nil {
		logger.Warn(orderDiscountLogTag, "GetDiscountByID Failed|Err:%v", err)
		return nil, err
	}

	return retInfo, nil
}

func (odm *OrderDiscountModel) GetDiscountList() ([]*OrderDiscount, error) {
	discountList, err := utils.SqlQuery(odm.sqlCli, orderDiscountTable, &OrderDiscount{}, "")
	if err != nil {
		logger.Warn(orderDiscountLogTag, "GetDiscountList Failed|Err:%v", err)
		return nil, err
	}

	return discountList.([]*OrderDiscount), nil
}
