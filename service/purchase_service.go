package service

import (
	"database/sql"
	"fmt"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/utils"
)

const (
	purchaseServiceLogTag = "PurchaseService"
)

type PurchaseService struct {
	sqlCli              *sql.DB
	wxUserModel         *model.WxUserModel
	supplierModel       *model.SupplierModel
	purchaseOrderModel  *model.PurchaseOrderModel
	purchaseDetailModel *model.PurchaseDetailModel
	goodsModel          *model.GoodsModel

	menuTypeMap map[uint32]*model.MenuType
}

func NewPurchaseService(sqlCli *sql.DB) *PurchaseService {
	supplierModel := model.NewSupplierModelWithDB(sqlCli)
	purchaseOrderModel := model.NewPurchaseOrderModelWithDB(sqlCli)
	purchaseDetailModel := model.NewPurchaseDetailModelWithDB(sqlCli)
	wxUserModel := model.NewWxUserModelWithDB(sqlCli)
	goodsModel := model.NewGoodsModelWithDB(sqlCli)
	return &PurchaseService{
		sqlCli:              sqlCli,
		supplierModel:       supplierModel,
		purchaseOrderModel:  purchaseOrderModel,
		purchaseDetailModel: purchaseDetailModel,
		wxUserModel:         wxUserModel,
		goodsModel:          goodsModel,
	}
}

func (ps *PurchaseService) GetSupplierMap() (map[uint32]*model.Supplier, error) {
	supplierList, err := ps.supplierModel.GetSupplier(0, "", "", "")
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "GetSupplierMap Failed|Err:%v", err)
		return nil, err
	}
	supplierMap := make(map[uint32]*model.Supplier)
	for _, supplier := range supplierList {
		supplierMap[supplier.ID] = supplier
	}
	return supplierMap, nil
}

func (ps *PurchaseService) GetSupplierList(name, phoneNumber string) ([]*model.Supplier, int64, error) {
	supplierList, err := ps.supplierModel.GetSupplier(0, name, phoneNumber, "")
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
	suppliers, err := ps.supplierModel.GetSupplier(supplierID, "", "", "")
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "BindSupplier GetSupplier Failed|Err:%v", err)
		return err
	}
	if len(suppliers) == 0 {
		return fmt.Errorf("供应商未找到|ID:%v", supplierID)
	}

	wxUser, err := ps.wxUserModel.GetWxUserByOpenID(openID)
	if err != nil {
		logger.Warn(userServiceLogTag, "GetWxUserByOpenID Failed|Err:%v", err)
		return err
	}
	if wxUser == nil {
		logger.Warn(userServiceLogTag, "WxUser NotExist|OpenID:%v", openID)
		return fmt.Errorf("要绑定的用户不存在，请确认OpenID是否正确")
	}

	err = ps.supplierModel.UpdateOpenID(supplierID, openID)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "UpdateOpenID GetSupplier Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) RenewSupplier(supplierID uint32, endTime int64) error {
	suppliers, err := ps.supplierModel.GetSupplier(supplierID, "", "", "")
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

func (ps *PurchaseService) GetPurchaseList(status int8, uid, purchaseID uint32, startTime, endTime int64,
	page, pageSize int32) ([]*model.PurchaseOrder, int32, map[uint32][]*model.PurchaseDetail, error) {
	purchaseList, err := ps.purchaseOrderModel.GetPurchaseOrderList(purchaseID, status, 0, uid, startTime, endTime, page, pageSize)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "GetPurchaseOrder Failed|Err:%v", err)
		return nil, 0, nil, err
	}
	purchaseCount, err := ps.purchaseOrderModel.GetPurchaseOrderCount(purchaseID, status, 0, uid, startTime, endTime)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetPurchaseOrderCount Failed|Err:%v", err)
		return nil, 0, nil, err
	}

	if len(purchaseList) == 0 {
		return make([]*model.PurchaseOrder, 0), 0, make(map[uint32][]*model.PurchaseDetail), nil
	}

	orderIDList := make([]uint32, 0, len(purchaseList))
	for _, purchase := range purchaseList {
		orderIDList = append(orderIDList, purchase.ID)
	}
	details, err := ps.purchaseDetailModel.GetPurchaseDetailByOrderList(orderIDList, 0)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderDetailByOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}

	detailMap := make(map[uint32][]*model.PurchaseDetail, 0)
	for _, detail := range details {
		if _, ok := detailMap[detail.PurchaseID]; ok == false {
			detailMap[detail.PurchaseID] = make([]*model.PurchaseDetail, 0)
		}
		detailMap[detail.PurchaseID] = append(detailMap[detail.PurchaseID], detail)
	}
	return purchaseList, purchaseCount, detailMap, nil
}

func (ps *PurchaseService) ApplyPurchaseOrder(purchase *model.PurchaseOrder, details []*model.PurchaseDetail) error {
	supplier, err := ps.supplierModel.GetCurrentSupplier()
	if err != nil || supplier == nil {
		logger.Warn(purchaseServiceLogTag, "ApplyPurchaseOrder GetCurrentSupplier Failed|Err:%v", err)
		return fmt.Errorf("请续期供应商")
	}

	tx, err := ps.sqlCli.Begin()
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ApplyPurchaseOrder Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	totalAmount := 0.0
	for _, item := range details {
		totalAmount += item.Price * item.ExpectNumber
	}

	purchase.Supplier = supplier.ID
	purchase.TotalAmount = totalAmount
	err = ps.purchaseOrderModel.InsertWithTx(tx, purchase)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "Insert Purchase Failed|Err:%v", err)
		return err
	}

	for _, item := range details {
		item.PurchaseID = purchase.ID
	}
	err = ps.purchaseDetailModel.BatchInsertWithTx(tx, details)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "BatchInsert PurchaseDetail Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) ReviewPurchaseOrder(purchaseID uint32) error {
	purchase, err := ps.purchaseOrderModel.GetPurchaseOrder(purchaseID)
	if err != nil || purchase == nil {
		logger.Warn(purchaseServiceLogTag, "GetPurchaseOrder Failed|Err:%v", err)
		return fmt.Errorf("采购订单不存在或状态错误")
	}
	if purchase.Status != enum.PurchaseNew {
		logger.Warn(purchaseServiceLogTag, "ReviewPurchaseOrder Status Error|ID:%v|Status:%v",
			purchaseID, purchase.Status)
		return fmt.Errorf("采购订单状态错误")
	}

	dao := &model.PurchaseOrder{ID: purchaseID, Status: enum.PurchaseReviewed}
	err = ps.purchaseOrderModel.UpdatePurchaseStatus(dao)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ReviewPurchaseOrder Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) ConfirmPurchaseOrder(purchaseID uint32) error {
	purchase, err := ps.purchaseOrderModel.GetPurchaseOrder(purchaseID)
	if err != nil || purchase == nil {
		logger.Warn(purchaseServiceLogTag, "GetPurchaseOrder Failed|Err:%v", err)
		return fmt.Errorf("采购订单不存在")
	}
	if purchase.Status != enum.PurchaseReviewed {
		logger.Warn(purchaseServiceLogTag, "ConfirmPurchaseOrder Status Error|ID:%v|Status:%v",
			purchaseID, purchase.Status)
		return fmt.Errorf("采购订单未审核")
	}

	dao := &model.PurchaseOrder{ID: purchaseID, Status: enum.PurchaseAccept}
	err = ps.purchaseOrderModel.UpdatePurchaseStatus(dao)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ConfirmPurchaseOrder Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ps *PurchaseService) ReceivePurchaseOrder(purchaseID uint32, details []*model.PurchaseDetail) error {
	purchase, err := ps.purchaseOrderModel.GetPurchaseOrder(purchaseID)
	if err != nil || purchase == nil {
		logger.Warn(purchaseServiceLogTag, "ReceivePurchase GetOrder Failed|Err:%v", err)
		return fmt.Errorf("采购订单未找到|ID:%v", purchaseID)
	}

	payAmount := 0.0
	goodsList := make([]*model.Goods, 0, len(details))
	for _, item := range details {
		payAmount += item.Price * item.ReceiveNumber
		goodsList = append(goodsList, &model.Goods{ID: item.GoodsID, Quantity: item.ReceiveNumber})
	}

	tx, err := ps.sqlCli.Begin()
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ApplyPurchaseOrder Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	err = ps.purchaseDetailModel.BatchUpdateDetailWithTx(tx, details)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ReceivePurchase BatchUpdateDetail Failed|Err:%v", err)
		return err
	}

	purchase.Status = enum.PurchaseReceived
	purchase.PayAmount = payAmount
	err = ps.purchaseOrderModel.UpdatePurchaseWithTx(tx, purchase, "status", "pay_amount")
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "ReceivePurchaseOrder Failed|Err:%v", err)
		return err
	}

	logger.Info(purchaseServiceLogTag, "BatchAddQuantity|List:%#v", goodsList)
	err = ps.goodsModel.BatchAddQuantityWithTx(tx, goodsList)
	if err != nil {
		logger.Warn(purchaseServiceLogTag, "BatchAddQuantityWithTx Failed|Err:%v", err)
		return err
	}

	return nil
}
