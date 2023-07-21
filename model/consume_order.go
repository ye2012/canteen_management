package model

import (
	"database/sql"
	"time"
)

type ApplyOrder struct {
	Id       int64     `json:"id"`
	Status   uint8     `json:"status"`
	Creator  uint32    `json:"creators"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type ApplyOrderModel struct {
	sqlCli *sql.DB
}

func NewApplyOrderModelWithDB(sqlCli *sql.DB) *ApplyOrderModel {
	return &ApplyOrderModel{
		sqlCli: sqlCli,
	}
}
