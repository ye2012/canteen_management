package server

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/service"
	"github.com/canteen_management/utils"
	"github.com/gin-gonic/gin"
)

const (
	purchaseServerLogTag = "PurchaseServer"
)

type PurchaseServer struct {
	storeService    *service.StoreService
	purchaseService *service.PurchaseService
}

func NewPurchaseServer(dbConf utils.Config) (*PurchaseServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(purchaseServerLogTag, "NewMysqlClient Failed|Err:%v", err)
		return nil, err
	}
	purchaseService := service.NewPurchaseService(sqlCli)
	storeService := service.NewStoreService(sqlCli)
	return &PurchaseServer{
		purchaseService: purchaseService,
		storeService:    storeService,
	}, nil
}

func (ps *PurchaseServer) RequestSupplierList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.SupplierListReq)

	supplierList, lastValidityTime, err := ps.purchaseService.GetSupplierList(req.Name, req.PhoneNumber)
	if err != nil {
		logger.Warn(purchaseServerLogTag, "GetSupplierList Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	retData := &dto.SupplierListRes{
		SupplierList:     ConvertToSupplier(supplierList),
		LastValidityTime: lastValidityTime,
	}
	res.Data = retData
}

func (ps *PurchaseServer) RequestModifySupplier(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifySupplierReq)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ps.purchaseService.AddSupplier(ConvertFromSupplierInfo(req.Supplier))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ps.purchaseService.UpdateSupplier(ConvertFromSupplierInfo(req.Supplier))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(purchaseServerLogTag, "RequestModifySupplier Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ps *PurchaseServer) RequestBindSupplier(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.BindSupplierReq)

	err := ps.purchaseService.BindSupplier(req.SupplierID, req.OpenID)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ps *PurchaseServer) RequestRenewSupplier(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.RenewSupplierReq)

	err := ps.purchaseService.RenewSupplier(req.SupplierID, req.EndTime)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ps *PurchaseServer) RequestApplyPurchase(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ApplyPurchaseReq)
	AdminUid := uint32(0)

	goodsMap, err := ps.storeService.GetGoodsMap()
	if err != nil {
		res.Code = enum.SystemError
		res.Msg = err.Error()
		return
	}

	details := ConvertFromApplyPurchase(req.GoodsList, goodsMap)
	purchaseOrder := &model.PurchaseOrder{Creator: AdminUid, Status: enum.PurchaseNew}
	err = ps.purchaseService.ApplyPurchaseOrder(purchaseOrder, details)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ps *PurchaseServer) RequestConfirmPurchase(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ConfirmPurchaseReq)

	err := ps.purchaseService.ConfirmPurchaseOrder(req.PurchaseID)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ps *PurchaseServer) RequestReceivePurchase(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ReceivePurchaseReq)
	goodsMap, err := ps.storeService.GetGoodsMap()
	if err != nil {
		res.Code = enum.SystemError
		res.Msg = err.Error()
		return
	}

	details := ConvertFromApplyPurchase(req.GoodsList, goodsMap)
	err = ps.purchaseService.ReceivePurchaseOrder(req.PurchaseID, details)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	ps.storeService.GetStoreTypeList()
}

func ConvertToSupplier(suppliers []*model.Supplier) []*dto.SupplierInfo {
	retList := make([]*dto.SupplierInfo, 0)
	for _, supplier := range suppliers {
		retInfo := &dto.SupplierInfo{
			SupplierID:       supplier.ID,
			Name:             supplier.Name,
			PhoneNumber:      supplier.PhoneNumber,
			IDNumber:         supplier.IDNumber,
			Location:         supplier.Location,
			ValidityDeadline: supplier.ValidityDeadline.Unix(),
			OpenID:           supplier.OpenID,
		}
		retList = append(retList, retInfo)
	}
	return retList
}

func ConvertFromSupplierInfo(supplierInfo *dto.SupplierInfo) *model.Supplier {
	return &model.Supplier{ID: supplierInfo.SupplierID, Name: supplierInfo.Name, PhoneNumber: supplierInfo.PhoneNumber,
		IDNumber: supplierInfo.IDNumber, Location: supplierInfo.Location}
}

func ConvertFromApplyPurchase(goodsList []*dto.PurchaseGoodsInfo, goodsMap map[uint32]*model.Goods) []*model.PurchaseDetail {
	detailList := make([]*model.PurchaseDetail, 0, len(goodsList))
	for _, purchaseGoods := range goodsList {
		goods, ok := goodsMap[purchaseGoods.GoodsID]
		if ok == false {
			continue
		}
		detail := &model.PurchaseDetail{
			GoodsID:       goods.ID,
			GoodsType:     goods.GoodsTypeID,
			ExpectAmount:  purchaseGoods.ExpectAmount,
			ReceiveAmount: 0,
			Price:         goods.Price,
		}
		detailList = append(detailList, detail)
	}
	return detailList
}
