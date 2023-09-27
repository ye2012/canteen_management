package model

import (
	"database/sql"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
	"time"
)

const (
	wxUserTable = "wx_user"

	wxUserModelLogTag = "WxUserModel"
)

type WxUser struct {
	ID                uint32    `json:"id"`
	OpenID            string    `json:"open_id"`
	PhoneNumber       string    `json:"phone_number"`
	OrderDiscountType uint8     `json:"order_discount_type"`
	CreateAt          time.Time `json:"created_at"`
	UpdateAt          time.Time `json:"updated_at"`
}

type WxUserModel struct {
	sqlCli *sql.DB
}

func NewWxUserModelWithDB(sqlCli *sql.DB) *WxUserModel {
	return &WxUserModel{
		sqlCli: sqlCli,
	}
}

func (wum *WxUserModel) Insert(dao *WxUser) error {
	id, err := utils.SqlInsert(wum.sqlCli, wxUserTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(wxUserModelLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (wum *WxUserModel) GetWxUserByOpenID(openID string) (*WxUser, error) {
	condition := " WHERE `open_id`=? "
	users, err := wum.GetWxUserByCondition(condition, openID)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (wum *WxUserModel) GetWxUserByCondition(condition string, params ...interface{}) ([]*WxUser, error) {
	retList, err := utils.SqlQuery(wum.sqlCli, wxUserTable, &WxUser{}, condition, params...)
	if err != nil {
		logger.Warn(wxUserModelLogTag, "GetUserByCondition Failed|Condition:%v|Params:%#v|Err:%v",
			condition, params, err)
		return nil, err
	}

	return retList.([]*WxUser), nil
}

func (wum *WxUserModel) UpdateWithTx(tx *sql.Tx, wxUser *WxUser, conditionTag string, params ...string) error {
	if tx == nil {
		err := utils.SqlUpdateWithUpdateTags(wum.sqlCli, wxUserTable, wxUser, conditionTag, params...)
		if err != nil {
			logger.Warn(wxUserModelLogTag, "SqlUpdateWithUpdateTags Failed|Err:%v", err)
			return err
		}
	} else {
		err := utils.SqlUpdateWithUpdateTags(tx, wxUserTable, wxUser, conditionTag, params...)
		if err != nil {
			logger.Warn(wxUserModelLogTag, "SqlUpdateWithUpdateTags Failed|Err:%v", err)
			return err
		}
	}
	return nil
}

func (wum *WxUserModel) Update(wxUser *WxUser, conditionTag string, params ...string) error {
	return wum.UpdateWithTx(nil, wxUser, conditionTag, params...)
}
