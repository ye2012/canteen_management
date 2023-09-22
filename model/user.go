package model

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	userTable = "user"

	userModelLogTag = "UserModel"
)

type User struct {
	ID                uint32    `json:"id"`
	UserName          string    `json:"user_name"`
	Password          string    `json:"password"`
	UnionID           string    `json:"union_id"`
	PhoneNumber       string    `json:"phone_number"`
	OrderDiscountType uint32    `json:"order_discount_type"`
	Role              uint32    `json:"role"`
	CreateAt          time.Time `json:"created_at"`
	UpdateAt          time.Time `json:"updated_at"`
}

type UserModel struct {
	sqlCli *sql.DB
}

func NewUserModelWithDB(sqlCli *sql.DB) *UserModel {
	return &UserModel{
		sqlCli: sqlCli,
	}
}

func (um *UserModel) Insert(dao *User) error {
	id, err := utils.SqlInsert(um.sqlCli, userTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(userModelLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (um *UserModel) GetUserByCondition(condition string, params ...interface{}) ([]*User, error) {
	retList, err := utils.SqlQuery(um.sqlCli, userTable, &User{}, condition, params...)
	if err != nil {
		logger.Warn(userModelLogTag, "GetUserByCondition Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*User), nil
}
