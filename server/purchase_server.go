package server

import (
	"github.com/canteen_management/conv"
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
	cartService     *service.CartService
}

func NewPurchaseServer(dbConf utils.Config) (*PurchaseServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(purchaseServerLogTag, "NewMysqlClient Failed|Err:%v", err)
		return nil, err
	}
	purchaseService := service.NewPurchaseService(sqlCli)
	storeService := service.NewStoreService(sqlCli)
	cartService := service.NewCartService(sqlCli)
	return &PurchaseServer{
		purchaseService: purchaseService,
		storeService:    storeService,
		cartService:     cartService,
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
		SupplierList:     conv.ConvertToSupplier(supplierList),
		LastValidityTime: lastValidityTime,
	}
	res.Data = retData
}

func (ps *PurchaseServer) RequestModifySupplier(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifySupplierReq)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ps.purchaseService.AddSupplier(conv.ConvertFromSupplierInfo(req.Supplier))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ps.purchaseService.UpdateSupplier(conv.ConvertFromSupplierInfo(req.Supplier))
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

func (ps *PurchaseServer) RequestPurchaseList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.PurchaseListReq)
	goodsMap, err := ps.storeService.GetGoodsMap()
	if err != nil {
		logger.Warn(purchaseServerLogTag, "GetGoodsMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}
	supplierMap, err := ps.purchaseService.GetSupplierMap()
	if err != nil {
		logger.Warn(purchaseServerLogTag, "GetSupplierMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	purchaseList, totalNumber, detailMap, err := ps.purchaseService.GetPurchaseList(req.Status, req.Uid, req.PurchaseID,
		req.StartTime, req.EndTime, req.Page, req.PageSize)
	if err != nil {
		logger.Warn(purchaseServerLogTag, "GetPurchaseList Failed|Err:%v", err)
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	logger.Debug(purchaseServerLogTag, "PurchaseListLen:%v", len(purchaseList))
	orderInfoList := conv.ConvertToPurchaseInfoList(purchaseList, detailMap, goodsMap, supplierMap)
	resData := &dto.PurchaseListRes{
		PurchaseList: orderInfoList,
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: totalNumber,
		},
	}
	res.Data = resData
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

	details := conv.ConvertFromApplyPurchase(req.GoodsList, goodsMap)
	purchaseOrder := &model.PurchaseOrder{Creator: AdminUid, Status: enum.PurchaseNew}
	err = ps.purchaseService.ApplyPurchaseOrder(purchaseOrder, details)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	ps.cartService.ClearCart(AdminUid, enum.CartTypePurchase)
}

func (ps *PurchaseServer) RequestReviewPurchase(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ReviewPurchaseReq)

	err := ps.purchaseService.ReviewPurchaseOrder(req.PurchaseID)
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

	details := conv.ConvertFromApplyPurchase(req.GoodsList, goodsMap)
	err = ps.purchaseService.ReceivePurchaseOrder(req.PurchaseID, details)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	ps.storeService.GetStoreTypeList()
}

func (ps *PurchaseServer) RequestApplyOutbound(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
}

func (ps *PurchaseServer) RequestOutboundOrderList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
}
