package model

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	tokenModelLogTag = "TokenModel"
	tokenTableName   = "token"

	md5TokenString = "canteen@2023"
)

type TokenModel struct {
	mysql *sql.DB
}

func NewTokenModelWithDB(mysql *sql.DB) *TokenModel {
	return &TokenModel{
		mysql: mysql,
	}
}

type TokenDAO struct {
	ID       uint32    `json:"id"`
	UID      uint32    `json:"uid"`
	AdminUid uint32    `json:"admin_uid"`
	Role     uint32    `json:"role"`
	Token    string    `json:"token"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

func (tm *TokenModel) RemoveUser(uid uint32) error {
	sqlStr := fmt.Sprintf("DELETE FROM `%v` WHERE `uid`=?", tokenTableName)

	_, err := tm.mysql.Exec(sqlStr, uid)
	if err != nil {
		logger.Warn(tokenModelLogTag, "delete token failed, err: %v", err)
	}
	return err
}

func (tm *TokenModel) RemoveAdminUser(adminUid uint32) error {
	sqlStr := fmt.Sprintf("DELETE FROM `%v` WHERE `admin_uid`=?", tokenTableName)

	_, err := tm.mysql.Exec(sqlStr, adminUid)
	if err != nil {
		logger.Warn(tokenModelLogTag, "delete token failed, err: %v", err)
	}
	return err
}

func (tm *TokenModel) Replace(dao *TokenDAO) error {
	id, err := utils.SqlReplace(tm.mysql, tokenTableName, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(tokenModelLogTag, "get new token ID failed, err: %v", err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (tm *TokenModel) Get(token string) *TokenDAO {
	ret := &TokenDAO{}
	err := utils.SqlQueryRow(tm.mysql, tokenTableName, ret, " WHERE `token` = ? ", token)
	if err != nil {
		logger.Warn(tokenModelLogTag, "get token failed: %v", err)
		return nil
	}
	return ret
}

func (tm *TokenModel) LoginSuccess(uid, adminUid, role uint32) *TokenDAO {
	nowTime := time.Now()
	token := fmt.Sprintf("%v%v%v%v", md5TokenString, uid, nowTime.UnixNano()/1e6, rand.Float64())
	token = utils.GetMD5Hex(token)

	dao := &TokenDAO{
		UID:      uid,
		AdminUid: adminUid,
		Role:     role,
		Token:    token,
	}
	tm.Replace(dao)
	return dao
}

func (tm *TokenModel) UserLoginSuccess(uid, role uint32) *TokenDAO {
	return tm.LoginSuccess(uid, 0, role)
}

func (tm *TokenModel) AdminLoginSuccess(adminUid, role uint32) *TokenDAO {
	return tm.LoginSuccess(0, adminUid, role)
}

func (tm *TokenModel) UpdateTokenWithTx(tx *sql.Tx, dao *TokenDAO, updateTags ...string) (err error) {
	if tx != nil {
		err = utils.SqlUpdateWithUpdateTags(tx, tokenTableName, dao, "id", updateTags...)
	} else {
		err = utils.SqlUpdateWithUpdateTags(tm.mysql, tokenTableName, dao, "id", updateTags...)
	}
	if err != nil {
		logger.Warn(tokenModelLogTag, "UpdateToken Failed|Err:%v", err)
		return err
	}
	return nil
}
