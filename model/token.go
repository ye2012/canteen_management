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

	tokenUID      = "`uid`"
	tokenToken    = "`token`"
	tokenPlatform = "`platform`"
	tokenCTime    = "`ctime`"
	tokenMTime    = "`mtime`"
	tokenID       = "`id`"
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
	ID       int64
	UID      int64
	Token    string
	Platform int
	CTime    time.Time
	MTime    time.Time
}

func (tm *TokenModel) Remove(uid int64, platform int) error {
	sqlStr := fmt.Sprintf("DELETE FROM %v WHERE %v=? AND %v=?", tokenTableName, tokenUID, tokenPlatform)

	_, err := tm.mysql.Exec(sqlStr, uid, platform)
	if err != nil {
		logger.Warn(tokenModelLogTag, "delete token failed, err: %v", err)
	}
	return err
}

func (tm *TokenModel) RemoveOtherPlatform(uid int64, platform int) error {
	sqlStr := fmt.Sprintf("DELETE FROM %v WHERE %v=? AND %v!=?", tokenTableName, tokenUID, tokenPlatform)

	_, err := tm.mysql.Exec(sqlStr, uid, platform)
	if err != nil {
		logger.Warn(tokenModelLogTag, "delete token failed, err: %v", err)
	}
	return err
}

func (tm *TokenModel) RemoveUser(uid int64) error {
	sqlStr := fmt.Sprintf("DELETE FROM %v WHERE %v=?", tokenTableName, tokenUID)

	_, err := tm.mysql.Exec(sqlStr, uid)
	if err != nil {
		logger.Warn(tokenModelLogTag, "delete token failed, err: %v", err)
	}
	return err
}

func (tm *TokenModel) Add(dao *TokenDAO) error {
	sqlStr := fmt.Sprintf("INSERT INTO %v (%v, %v, %v, %v, %v) VALUES(?, ?, ?, ?, ?)", tokenTableName,
		tokenUID, tokenToken, tokenPlatform, tokenCTime, tokenMTime)
	res, err := tm.mysql.Exec(sqlStr, dao.UID, dao.Token, dao.Platform, dao.CTime, dao.MTime)
	if err != nil {
		logger.Warn(tokenModelLogTag, "insert token failed, err: %v", err)
		return err
	}
	dao.ID, err = res.LastInsertId()
	if err != nil {
		logger.Warn(tokenModelLogTag, "get new token ID failed, err: %v", err)
		return err
	}
	return nil
}

func (tm *TokenModel) Replace(dao *TokenDAO) error {
	sqlStr := fmt.Sprintf("REPLACE INTO %v (%v, %v, %v, %v, %v) VALUES(?, ?, ?, ?, ?)", tokenTableName,
		tokenUID, tokenToken, tokenPlatform, tokenCTime, tokenMTime)
	res, err := tm.mysql.Exec(sqlStr, dao.UID, dao.Token, dao.Platform, dao.CTime, dao.MTime)
	if err != nil {
		logger.Warn(tokenModelLogTag, "replace token failed, err: %v", err)
		return err
	}
	dao.ID, err = res.LastInsertId()
	if err != nil {
		logger.Warn(tokenModelLogTag, "get new token ID failed, err: %v", err)
		return err
	}
	return nil
}

func (tm *TokenModel) Get(token string) *TokenDAO {
	sqlStr := fmt.Sprintf("SELECT %v, %v, %v, %v, %v, %v FROM %v WHERE %v=?",
		tokenID, tokenUID, tokenToken, tokenPlatform, tokenCTime, tokenMTime,
		tokenTableName, tokenToken)
	row := tm.mysql.QueryRow(sqlStr, token)
	ret := &TokenDAO{}
	err := row.Scan(&ret.ID, &ret.UID, &ret.Token, &ret.Platform, &ret.CTime, &ret.MTime)
	if err != nil {
		logger.Warn(tokenModelLogTag, "get token failed: %v", err)
		return nil
	}
	return ret
}

func (tm *TokenModel) GetByUidPlatform(uid int64, platform int) *TokenDAO {
	sqlStr := fmt.Sprintf("SELECT %v, %v, %v, %v, %v, %v FROM %v WHERE %v=? AND %v = ?",
		tokenID, tokenUID, tokenToken, tokenPlatform, tokenCTime, tokenMTime,
		tokenTableName, tokenUID, tokenPlatform)
	row := tm.mysql.QueryRow(sqlStr, uid, platform)
	ret := &TokenDAO{}
	err := row.Scan(&ret.ID, &ret.UID, &ret.Token, &ret.Platform, &ret.CTime, &ret.MTime)
	if err != nil {
		logger.Warn(tokenModelLogTag, "get token failed: %v", err)
		return nil
	}
	return ret
}

func (tm *TokenModel) UpdateModifyTime(token *TokenDAO) {
	sqlStr := fmt.Sprintf("UPDATE %v SET %v=? WHERE %v=?",
		tokenTableName, tokenMTime, tokenID)
	_, err := tm.mysql.Exec(sqlStr, token.MTime, token.ID)
	if err != nil {
		logger.Warn(tokenModelLogTag, "UpdateModifyTime failed: %v", err)
	}
}

func (tm *TokenModel) LoginSuccess(uid int64, platform int) *TokenDAO {
	nowTime := time.Now()
	token := fmt.Sprintf("%v%v%v%v", md5TokenString, uid, nowTime.UnixNano()/1e6, rand.Float64())
	token = utils.GetMD5Hex(token)

	dao := &TokenDAO{
		UID:      uid,
		Token:    token,
		Platform: platform,
		CTime:    nowTime,
		MTime:    nowTime,
	}

	tm.RemoveOtherPlatform(uid, platform)
	// 添加token到数据库
	tm.Replace(dao)

	return dao
}
