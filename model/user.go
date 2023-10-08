package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	adminUserTable = "admin_user"

	adminUserModelLogTag = "AdminUserModel"
)

var (
	adminUserUpdateTags = []string{"nick_name", "phone_number", "role"}
)

type AdminUser struct {
	ID          uint32    `json:"id"`
	UserName    string    `json:"user_name"`
	Password    string    `json:"password"`
	NickName    string    `json:"nick_name"`
	OpenID      string    `json:"open_id"`
	PhoneNumber string    `json:"phone_number"`
	Role        uint32    `json:"role"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

type AdminUserModel struct {
	sqlCli *sql.DB
}

func NewAdminUserModelWithDB(sqlCli *sql.DB) *AdminUserModel {
	return &AdminUserModel{
		sqlCli: sqlCli,
	}
}

func (aum *AdminUserModel) Insert(dao *AdminUser) error {
	id, err := utils.SqlInsert(aum.sqlCli, adminUserTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(adminUserModelLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (aum *AdminUserModel) DeleteByID(id uint32) error {
	sqlStr := fmt.Sprintf("DELETE FROM `%v` WHERE `id` = ?", adminUserTable)
	_, err := aum.sqlCli.Exec(sqlStr, id)
	if err != nil {
		logger.Warn(adminUserModelLogTag, "DeleteByID Failed|ID:%v|Err:%v", id, err)
		return err
	}

	return nil
}

func (aum *AdminUserModel) GetAdminUserByCondition(condition string, params ...interface{}) ([]*AdminUser, error) {
	retList, err := utils.SqlQuery(aum.sqlCli, adminUserTable, &AdminUser{}, condition, params...)
	if err != nil {
		logger.Warn(adminUserModelLogTag, "GetAdminUserByCondition Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*AdminUser), nil
}

func (aum *AdminUserModel) UpdateAdminUserByCondition(adminUser *AdminUser, conditionTag string, updateTags ...string) error {
	err := utils.SqlUpdateWithUpdateTags(aum.sqlCli, adminUserTable, adminUser, conditionTag, updateTags...)
	if err != nil {
		logger.Warn(adminUserModelLogTag, "UpdateAdminUserByCondition Failed|Err:%v", err)
		return err
	}

	return nil
}

func (aum *AdminUserModel) UpdateAdminUserInfo(adminUser *AdminUser, conditionTag string) error {
	return aum.UpdateAdminUserByCondition(adminUser, conditionTag, adminUserUpdateTags...)
}
