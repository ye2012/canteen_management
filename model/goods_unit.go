package model

import (
	"database/sql"
	"time"
)

type GoodsUnit struct {
	Id            int64     `json:"id"`
	GoodsUnitName string    `json:"goods_unit_name"`
	IsBaseUnit    bool      `json:"is_base_unit"`
	CreateAt      time.Time `json:"created_at"`
	UpdateAt      time.Time `json:"updated_at"`
}

type GoodsUnitModel struct {
	sqlCli *sql.DB
}

func NewGoodsUnitModelWithDB(sqlCli *sql.DB) *GoodsUnitModel {
	return &GoodsUnitModel{
		sqlCli: sqlCli,
	}
}
