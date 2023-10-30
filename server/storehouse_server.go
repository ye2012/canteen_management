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
	storeServerLogTag = "StoreServer"
)

type StorehouseServer struct {
	storeService     *service.StoreService
	inventoryService *service.InventoryService
	userService      *service.UserService
}

func NewStorehouseServer(dbConf utils.Config) (*StorehouseServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(storeServerLogTag, "NewStorehouseServer Failed|Err:%v", err)
		return nil, err
	}
	storeService := service.NewStoreService(sqlCli)
	inventoryService := service.NewInventoryService(sqlCli)
	userService := service.NewUserService(sqlCli)
	return &StorehouseServer{
		storeService:     storeService,
		inventoryService: inventoryService,
		userService:      userService,
	}, nil
}

func (ss *StorehouseServer) RequestGoodsTypeList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GoodsTypeListReq)
	list, err := ss.storeService.GetGoodsTypeList()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	skip := (req.Page - 1) * req.PageSize
	if skip < int32(len(list)) {
		list = list[skip:]
	} else {
		list = make([]*model.GoodsType, 0)
	}
	if int32(len(list)) > req.PageSize {
		list = list[:req.PageSize]
	}

	res.Data = &dto.GoodsTypeListRes{
		GoodsTypeList: conv.ConvertToGoodsTypeInfoList(list),
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: int32(len(list)),
		},
	}
}

func (ss *StorehouseServer) RequestModifyGoodsType(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyGoodsTypeReq)
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ss.storeService.AddGoodsType(conv.ConvertFromGoodsTypeInfo(req.GoodsType))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ss.storeService.UpdateGoodsType(conv.ConvertFromGoodsTypeInfo(req.GoodsType))
		if err != nil {
			res.Code = enum.SqlError
			res.Msg = err.Error()
			return
		}
	case enum.OperateTypeDel:
		err := ss.storeService.DeleteGoodsType(req.GoodsType.GoodsTypeID)
		if err != nil {
			res.Code = enum.SqlError
			res.Msg = err.Error()
			return
		}
	default:
		logger.Warn(storeServerLogTag, "RequestModifyGoodsType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ss *StorehouseServer) RequestGoodsList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GoodsListReq)
	goodsList, goodsCount, err := ss.storeService.GoodsList(req.GoodsTypeID, req.StoreTypeID, req.Page, req.PageSize)
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.GoodsListRes{
		GoodsList: conv.ConvertToGoodsInfoList(goodsList),
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: goodsCount,
		},
	}
}

func (ss *StorehouseServer) RequestModifyGoods(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyGoodsInfoReq)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ss.storeService.AddGoods(conv.ConvertFromGoodsInfo(req.Goods))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ss.storeService.UpdateGoods(conv.ConvertFromGoodsInfo(req.Goods))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeDel:
		err := ss.storeService.DeleteGoods(req.Goods.GoodsID)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(storeServerLogTag, "RequestModifyStoreType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ss *StorehouseServer) RequestGoodsPriceList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GoodsPriceListReq)
	goodsTypeMap, err := ss.storeService.GetGoodsTypeMap()
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	goodsList, goodsCount, err := ss.storeService.GoodsList(req.GoodsTypeID, req.StoreTypeID, req.Page, req.PageSize)
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.GoodsPriceListRes{
		GoodsPriceList: conv.ConvertToGoodsPriceList(goodsList, goodsTypeMap),
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: goodsCount,
		},
	}
}

func (ss *StorehouseServer) RequestModifyGoodsPrice(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyGoodsPriceReq)

	err := ss.storeService.UpdateGoodsPrice(req.GoodsID, req.PriceMap)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ss *StorehouseServer) RequestStoreTypeList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	list, err := ss.storeService.GetStoreTypeList()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.StoreTypeListRes{
		StoreTypeList: conv.ConvertToStoreTypeInfoList(list),
	}
}

func (ss *StorehouseServer) RequestModifyStoreType(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyStoreTypeReq)
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ss.storeService.AddStoreType(conv.ConvertFromStoreTypeInfo(req.StoreType))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ss.storeService.UpdateStoreType(conv.ConvertFromStoreTypeInfo(req.StoreType))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeDel:
		err := ss.storeService.DeleteStoreType(req.StoreType.StoreTypeID)
		if err != nil {
			res.Code = enum.SqlError
			res.Msg = err.Error()
			return
		}
	default:
		logger.Warn(storeServerLogTag, "RequestModifyStoreType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ss *StorehouseServer) RequestGoodsNodeList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GoodsNodeListReq)
	uid := req.Uid
	goodsMap, err := ss.storeService.GetGoodsMap()
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
	goodsTypeList, err := ss.storeService.GetGoodsTypeList()
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	goodsSelectedMap, totalCost, totalGoods, cartID := make(map[string]float64), 0.0, 0.0, uint32(0)
	if uid != 0 {
		cart, cartDetails, err := ss.storeService.GetCart(uid, req.CartType)
		if err != nil {
			res.Code = enum.SystemError
			return
		}
		if cart != nil {
			cartID = cart.ID
		}
		for _, detail := range cartDetails {
			goodsSelectedMap[detail.ItemID] = detail.Quantity
			goodsID, _ := conv.ConvertGoodsID(detail.ItemID)
			goods, ok := goodsMap[goodsID]
			if ok {
				totalCost += goods.Price * detail.Quantity
			}
			if detail.Quantity > 0 {
				totalGoods += 1
			}
		}
	}

	retData := &dto.GoodsNodeListRes{
		GoodsList:  conv.ConvertGoodsListToGoodsNode(goodsMap, goodsTypeList, goodsSelectedMap),
		GoodsMap:   goodsSelectedMap,
		TotalGoods: totalGoods,
		TotalCost:  totalCost,
		CartID:     cartID,
	}
	res.Data = retData
}

func (ss *StorehouseServer) RequestInventoryOrderList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.InventoryListReq)
	goodsMap, err := ss.storeService.GetGoodsMap()
	if err != nil {
		logger.Warn(storeServerLogTag, "GetGoodsMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}
	adminMap, err := ss.userService.GetAdminMap()
	if err != nil {
		logger.Warn(storeServerLogTag, "GetAdminMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	orderList, orderCount, detailMap, err := ss.inventoryService.InventoryOrderList(req.InventoryID, req.Uid,
		req.Status, req.StartTime, req.EndTime, req.Page, req.PageSize)
	if err != nil {
		logger.Warn(storeServerLogTag, "GetInventoryList Failed|Err:%v", err)
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	orderInfoList := conv.ConvertToInventoryInfoList(orderList, detailMap, goodsMap, adminMap)
	resData := &dto.InventoryListRes{
		InventoryList: orderInfoList,
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: orderCount,
		},
	}
	res.Data = resData
}

func (ss *StorehouseServer) RequestStartInventory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.StartInventoryReq)
	goodsMap, err := ss.storeService.GetGoodsMap()
	if err != nil {
		logger.Warn(storeServerLogTag, "GetGoodsMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}
	goodsTypeMap, err := ss.storeService.GetGoodsTypeList()
	if err != nil {
		logger.Warn(storeServerLogTag, "GetGoodsTypeList Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	retData := &dto.StartInventoryRes{}
	inventoryOrders, _, inventoryDetailMap, err := ss.inventoryService.InventoryOrderList(0,
		req.Uid, 0, 0, 0, 1, 10)
	if err != nil {
		logger.Warn(storeServerLogTag, "StartInventory GetOrder Failed|Err:%v", err)
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
	if len(inventoryOrders) > 0 {
		retData.InventoryOrder = conv.ConvertGoodsListToInventoryNode(inventoryOrders[0],
			inventoryDetailMap[inventoryOrders[0].ID], goodsMap, goodsTypeMap)
	}
	if req.New && len(inventoryOrders) == 0 {
		inventoryOrder, detailList, err := ss.inventoryService.StartInventory(req.Uid)
		if err != nil {
			logger.Warn(storeServerLogTag, "StartInventory Failed|Err:%v", err)
			res.Code = enum.SqlError
			res.Msg = err.Error()
			return
		}
		retData.InventoryOrder = conv.ConvertGoodsListToInventoryNode(inventoryOrder, detailList, goodsMap, goodsTypeMap)
	}

	res.Data = retData
}

func (ss *StorehouseServer) RequestInventory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.InventoryReq)

	inventoryDetail := conv.ConvertFromApplyInventory(req.InventoryGoodsInfo)
	err := ss.inventoryService.UpdateInventory(inventoryDetail)
	if err != nil {
		logger.Warn(storeServerLogTag, "UpdateInventory Failed|Err:%v", err)
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ss *StorehouseServer) RequestApplyInventory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ApplyInventoryReq)

	err := ss.inventoryService.ApplyInventory(req.InventoryID)
	if err != nil {
		logger.Warn(storeServerLogTag, "ApplyInventory Failed|Err:%v", err)
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ss *StorehouseServer) RequestConfirmInventory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ConfirmInventoryReq)

	err := ss.inventoryService.ConfirmInventory(req.InventoryID, req.Uid)
	if err != nil {
		logger.Warn(storeServerLogTag, "ConfirmInventory Failed|Err:%v", err)
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ss *StorehouseServer) RequestReviewInventory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ReviewInventoryReq)
	err := ss.inventoryService.ReviewInventory(req.InventoryID)
	if err != nil {
		logger.Warn(storeServerLogTag, "ReviewInventory Failed|Err:%v", err)
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (ss *StorehouseServer) RequestGoodsHistory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GoodsHistoryReq)
	history, count, err := ss.storeService.GetGoodsHistoryList(req.GoodsID, req.ChangeType, req.StartTime, req.EndTime,
		req.Page, req.PageSize)
	if err != nil {
		logger.Warn(storeServerLogTag, "RequestGoodsHistory Failed|Err:%v", err)
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	retData := &dto.GoodsHistoryRes{
		History: conv.ConvertToGoodsHistoryInfoList(history),
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: count,
		},
	}
	res.Data = retData
}
