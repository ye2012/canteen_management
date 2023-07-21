package model

import (
	"database/sql"
	"time"
)

type Goods struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	GoodsTypeID uint32    `json:"goods_type_id"`
	Picture     string    `json:"picture"`
	BatchSize   float64   `json:"batch_size"`
	BatchUnit   string    `json:"batch_unit_id"`
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
