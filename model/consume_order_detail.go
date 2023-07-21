package model

import "database/sql"

type ApplyOrderDetail struct {
	Id           int64   `json:"id"`
	ApplyOrderID uint32  `json:"apply_order_id"`
	GoodsID      uint32  `json:"goods_id"`
	ApplyAmount  float64 `json:"apply_amount"`
	SendAmount   float64 `json:"send_amount"`
}

type ApplyOrderDetailModel struct {
	sqlCli *sql.DB
}

func NewApplyOrderDetailModelWithDB(sqlCli *sql.DB) *ApplyOrderDetailModel {
	return &ApplyOrderDetailModel{
		sqlCli: sqlCli,
	}
}
