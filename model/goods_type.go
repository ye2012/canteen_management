package model

import (
	"database/sql"
	"time"
)

type GoodsType struct {
	Id            int64     `json:"id"`
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
