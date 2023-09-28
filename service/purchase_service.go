package service

import (
	"database/sql"
	"fmt"
	"github.com/canteen_management/enum"
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

func (ps *PurchaseService) GetSupplierList(name, phoneNumber string) ([]*model.Supplier, int64, error) {
	supplierList, err := ps.supplierModel.GetSupplier(0, name, phoneNumber)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "GetSupplier Failed|Err:%v", err)
		return nil, 0, err
	}
	lastTime, err := ps.supplierModel.GetLastValidityTime()
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "GetLastValidityTime Failed|Err:%v", err)
		return nil, 0, err
	}
	return supplierList, lastTime, nil
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

func (ps *PurchaseService) RenewSupplier(supplierID uint32, endTime int64) error {
	suppliers, err := ps.supplierModel.GetSupplier(supplierID, "", "")
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "RenewSupplier GetSupplier Failed|Err:%v", err)
		return err
	}
	if len(suppliers) == 0 {
		return fmt.Errorf("供应商未找到|ID:%v", supplierID)
	}

	lastValidityTime, err := ps.supplierModel.GetLastValidityTime()
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "GetLastValidityTime Failed|Err:%v", err)
		return err
	}
	if endTime < lastValidityTime {
		return fmt.Errorf("续期时间应该晚于当前供应商")
	}

	err = ps.supplierModel.UpdateValidityTime(supplierID, endTime)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "UpdateValidityTime GetSupplier Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) GetPurchaseOrder() {

}

func (ps *PurchaseService) ApplyPurchaseOrder(purchase *model.PurchaseOrder, details []*model.PurchaseDetail) error {
	supplier, err := ps.supplierModel.GetCurrentSupplier()
	if err != nil || supplier == nil {
		logger.Warn(purchaseServiceLogTag, "ApplyPurchaseOrder GetCurrentSupplier Failed|Err:%v", err)
		return fmt.Errorf("请续期供应商")
	}

	totalAmount := 0.0
	for _, item := range details {
		totalAmount += item.Price
	}

	purchase.Supplier = supplier.ID
	purchase.TotalAmount = totalAmount
	err = ps.purchaseOrderModel.Insert(purchase)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "Insert Purchase Failed|Err:%v", err)
		return err
	}

	for _, item := range details {
		item.PurchaseOrderID = purchase.ID
	}
	err = ps.purchaseDetailModel.BatchInsert(details)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "BatchInsert PurchaseDetail Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) ConfirmPurchaseOrder(purchaseID uint32) error {
	dao := &model.PurchaseOrder{ID: purchaseID, Status: enum.PurchaseAccept}
	err := ps.purchaseOrderModel.UpdatePurchaseStatus(dao)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ConfirmPurchaseOrder Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) ReceivePurchaseOrder(purchaseID uint32, details []*model.PurchaseDetail) error {
	purchases, err := ps.purchaseOrderModel.GetPurchaseOrder(purchaseID, enum.PurchaseOrderStatusAll, 0, 0, 0)
	if err != nil || len(purchases) == 0 {
		logger.Warn(purchaseServiceLogTag, "ReceivePurchase GetOrder Failed|Err:%v", err)
		return fmt.Errorf("采购订单未找到|ID:%v", purchaseID)
	}

	payAmount := 0.0
	for _, item := range details {
		payAmount += item.Price * item.ReceiveAmount
	}

	err = ps.purchaseDetailModel.BatchUpdateDetail(details)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ReceivePurchase BatchUpdateDetail Failed|Err:%v", err)
		return err
	}

	purchase := purchases[0]
	purchase.Status = enum.PurchaseReceived
	purchase.PayAmount = payAmount
	err = ps.purchaseOrderModel.UpdatePurchase(purchase, "status", "pay_amount")
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ReceivePurchaseOrder Failed|Err:%v", err)
		return err
	}
	return nil
}
