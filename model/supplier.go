package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/utils"
)

const (
	supplierTable  = "supplier"
	supplierLogTag = "SupplierModel"
)

var (
	supplierUpdateTags = []string{"name", "phone_number", "id_number", "location", "validity_deadline"}
)

type Supplier struct {
	ID               uint32    `json:"id"`
	Name             string    `json:"name"`
	PhoneNumber      string    `json:"phone_number"`
	IDNumber         string    `json:"id_number"`
	Location         string    `json:"location"`
	ValidityDeadline time.Time `json:"validity_deadline"`
	OpenID           string    `json:"open_id"`
	CreateAt         time.Time `json:"created_at"`
	UpdateAt         time.Time `json:"updated_at"`
}

type SupplierModel struct {
	sqlCli *sql.DB
}

func NewSupplierModelWithDB(sqlCli *sql.DB) *SupplierModel {
	return &SupplierModel{
		sqlCli: sqlCli,
	}
}

func (sm *SupplierModel) Insert(dao *Supplier) error {
	id, err := utils.SqlInsert(sm.sqlCli, supplierTable, dao, "id", "created_at", "updated_at")
	if err != nil {
		logger.Warn(supplierLogTag, "Insert Failed|Dao:%+v|Err:%v", dao, err)
		return err
	}
	dao.ID = uint32(id)
	return nil
}

func (sm *SupplierModel) GetCurrentSupplier() (*Supplier, error) {
	condition := " WHERE `validity_deadline` > CURRENT_TIMESTAMP() ORDER BY `validity_deadline` ASC LIMIT 1 "
	supplier := &Supplier{}
	err := utils.SqlQueryRow(sm.sqlCli, supplierTable, supplier, condition)
	if err != nil {
		logger.Warn(supplierLogTag, "GetCurrentSupplier Failed|Err:%v", err)
		return nil, err
	}

	return supplier, nil
}

func (sm *SupplierModel) GetSupplier(id uint32, name, phoneNumber, openID string) ([]*Supplier, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if id != 0 {
		condition += " AND `id` = ? "
		params = append(params, id)
	}
	if name != "" {
		condition += " AND `name` = ? "
		params = append(params, name)
	}
	if phoneNumber != "" {
		condition += " AND `phone_number` = ? "
		params = append(params, phoneNumber)
	}
	if openID != "" {
		condition += " AND `open_id` = ? "
		params = append(params, openID)
	}

	retList, err := utils.SqlQuery(sm.sqlCli, supplierTable, &Supplier{}, condition, params...)
	if err != nil {
		logger.Warn(supplierLogTag, "GetSupplier Failed|Err:%v", err)
		return nil, err
	}

	return retList.([]*Supplier), nil
}

func (sm *SupplierModel) UpdateSupplier(dao *Supplier) error {
	err := utils.SqlUpdateWithUpdateTags(sm.sqlCli, supplierTable, dao, "id", supplierUpdateTags...)
	if err != nil {
		logger.Warn(supplierLogTag, "UpdateSupplier Failed|Err:%v", err)
		return err
	}
	return nil
}

func (sm *SupplierModel) UpdateOpenID(id uint32, openID string) error {
	dao := &Supplier{ID: id, OpenID: openID}
	err := utils.SqlUpdateWithUpdateTags(sm.sqlCli, supplierTable, dao, "id", "open_id")
	if err != nil {
		logger.Warn(supplierLogTag, "UpdateOpenID Failed|Err:%v", err)
		return err
	}
	return nil
}

func (sm *SupplierModel) UpdateValidityTime(id uint32, endTime int64) error {
	dao := &Supplier{ID: id, ValidityDeadline: time.Unix(endTime, 0)}
	err := utils.SqlUpdateWithUpdateTags(sm.sqlCli, supplierTable, dao, "id", "validity_deadline")
	if err != nil {
		logger.Warn(supplierLogTag, "UpdateValidityTime Failed|Err:%v", err)
		return err
	}
	return nil
}

func (sm *SupplierModel) GetLastValidityTime() (int64, error) {
	sqlStr := fmt.Sprintf("SELECT MAX(`validity_deadline`) FROM %v ", supplierTable)
	row := sm.sqlCli.QueryRow(sqlStr)
	lastTime := time.Time{}
	err := row.Scan(&lastTime)
	if err != nil {
		logger.Warn(supplierLogTag, "Scan LastValidityTime Failed|Err:%v", err)
		return 0, err
	}
	return lastTime.Unix(), nil
}
