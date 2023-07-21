package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id          int64     `json:"id"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	Role        uint32    `json:"role"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

type UserModel struct {
	sqlCli *sql.DB
}

func NewUserModelWithDB(sqlCli *sql.DB) *UserModel {
	return &UserModel{
		sqlCli: sqlCli,
	}
}
