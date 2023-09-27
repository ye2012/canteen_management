package service

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	purchaseServiceLogTag = "PurchaseService"
)

type PurchaseService struct {
	supplierModel       *model.SupplierModel
	purchaseOrderModel  *model.PurchaseOrderModel
	purchaseDetailModel *model.PurchaseDetailModel

	menuTypeMap map[uint32]*model.MenuType
}

func NewPurchaseService(sqlCli *sql.DB) *PurchaseService {
	supplierModel := model.NewSupplierModelWithDB(sqlCli)
	purchaseOrderModel := model.NewPurchaseOrderModelWithDB(sqlCli)
	purchaseDetailModel := model.NewPurchaseDetailModelWithDB(sqlCli)
	return &PurchaseService{
		supplierModel:       supplierModel,
		purchaseOrderModel:  purchaseOrderModel,
		purchaseDetailModel: purchaseDetailModel,
	}
}

func (ps *PurchaseService) GetSupplierList(name, phoneNumber string) ([]*model.Supplier, error) {
	supplierList, err := ps.supplierModel.GetSupplier(0, name, phoneNumber)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "GetSupplier Failed|Err:%v", err)
		return nil, err
	}
	return supplierList, nil
}

func (ps *PurchaseService) AddSupplier(supplier *model.Supplier) error {
	err := ps.supplierModel.Insert(supplier)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "Insert Supplier Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) UpdateSupplier(supplier *model.Supplier) error {
	err := ps.supplierModel.UpdateSupplier(supplier)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "UpdateSupplier Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) BindSupplier(supplierID uint32, openID string) error {
	suppliers, err := ps.supplierModel.GetSupplier(supplierID, "", "")
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "BindSupplier GetSupplier Failed|Err:%v", err)
		return err
	}
	if len(suppliers) == 0 {
		return fmt.Errorf("供应商未找到|ID:%v", supplierID)
	}

	err = ps.supplierModel.UpdateOpenID(supplierID, openID)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "UpdateOpenID GetSupplier Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) GetPurchaseOrder() {

}
