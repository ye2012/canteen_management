package service

import (
	"database/sql"

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
	supplierList, err := ps.supplierModel.GetSupplier(name, phoneNumber)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "GetSupplier Failed|Err:%v", err)
		return nil, err
	}
	return supplierList, nil
}

func (ps *PurchaseService) AddSupplier() {

}

func (ps *PurchaseService) GetPurchaseOrder() {

}
