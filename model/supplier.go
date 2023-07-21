package model

import (
	"database/sql"
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

func (sm *SupplierModel) GetSupplier(name, phoneNumber string) ([]*Supplier, error) {
	var params []interface{}
	condition := " WHERE 1=1 "
	if name != "" {
		condition += " AND `name` = ? "
		params = append(params, name)
	}
	if phoneNumber != "" {
		condition += " AND `phone_number` = ? "
		params = append(params, phoneNumber)
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
