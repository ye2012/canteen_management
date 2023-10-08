package server

import (
	"github.com/canteen_management/conv"
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/service"
	"github.com/canteen_management/utils"
	"github.com/gin-gonic/gin"
)

const (
	storeServerLogTag = "StoreServer"
)

type StorehouseServer struct {
	storeService *service.StoreService
}

func NewStorehouseServer(dbConf utils.Config) (*StorehouseServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(storeServerLogTag, "NewStorehouseServer Failed|Err:%v", err)
		return nil, err
	}
	storeService := service.NewStoreService(sqlCli)
	return &StorehouseServer{
		storeService: storeService,
	}, nil
}

func (ss *StorehouseServer) RequestGoodsTypeList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	list, err := ss.storeService.GetGoodsTypeList()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.GoodsTypeListRes{
		GoodsTypeList: conv.ConvertToGoodsTypeInfoList(list),
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
	//ctx.Request.ParseMultipartForm(1024)
	//imgFile, _, err := ctx.Request.FormFile("img")
	//if err != nil {
	//	logger.Warn(storeServerLogTag, "ReadFrom File Failed|Err:%v", err)
	//	return
	//}
	//defer imgFile.Close()
	//
	//imgFilePath := config.Config.FileStorePath + "/Goods/" + req.Goods.GoodsName + ".jpg"
	//image, err := os.Create(imgFilePath)
	//if err != nil {
	//	logger.Warn(storeServerLogTag, "ReadFrom File Failed|Err:%v", err)
	//	return
	//}
	//defer image.Close()
	//
	//_, err = io.Copy(image, imgFile)
	//if err != nil {
	//	logger.Warn(storeServerLogTag, "Copy File Failed|Err:%v", err)
	//	return
	//}

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
	default:
		logger.Warn(storeServerLogTag, "RequestModifyStoreType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ss *StorehouseServer) RequestGoodsPriceList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GoodsPriceListReq)
	goodsList, goodsCount, err := ss.storeService.GoodsList(req.GoodsTypeID, req.StoreTypeID, req.Page, req.PageSize)
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.GoodsPriceListRes{
		GoodsPriceList: conv.ConvertToGoodsPriceList(goodsList),
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
	default:
		logger.Warn(storeServerLogTag, "RequestModifyStoreType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ss *StorehouseServer) RequestGoodsNodeList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GoodsNodeListReq)
	uid := uint32(1)
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

	goodsQuantityMap, totalCost, totalGoods := make(map[string]float64), 0.0, 0.0
	if uid != 0 {
		_, cartDetails, err := ss.storeService.GetCart(uid, req.CartType)
		if err != nil {
			res.Code = enum.SystemError
			return
		}
		for _, detail := range cartDetails {
			goodsQuantityMap[detail.ItemID] = detail.Quantity
			goodsID, _ := conv.ConvertGoodsID(detail.ItemID)
			goods, ok := goodsMap[goodsID]
			if ok {
				totalCost += goods.Price * detail.Quantity
			}
			totalGoods += detail.Quantity
		}
	}

	retData := &dto.GoodsNodeListRes{
		GoodsList:  conv.ConvertGoodsListToGoodsNode(goodsMap, goodsTypeList),
		GoodsMap:   goodsQuantityMap,
		TotalGoods: totalGoods,
		TotalCost:  totalCost,
	}
	res.Data = retData
}

func (ss *StorehouseServer) RequestApplyOutbound(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
}

func (ss *StorehouseServer) RequestOutboundOrderList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
}

func (ss *StorehouseServer) RequestInventoryOrder(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
}

func (ss *StorehouseServer) RequestStartInventory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
}

func (ss *StorehouseServer) RequestApplyInventory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
}

func (ss *StorehouseServer) RequestConfirmInventory(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
}
