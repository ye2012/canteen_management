package server

import (
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
		GoodsTypeList: ConvertToGoodsTypeInfoList(list),
	}
}

func (ss *StorehouseServer) RequestModifyGoodsType(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyGoodsTypeReq)
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ss.storeService.AddGoodsType(ConvertFromGoodsTypeInfo(req.GoodsType))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ss.storeService.UpdateGoodsType(ConvertFromGoodsTypeInfo(req.GoodsType))
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

	extraPage := int32(1)
	if goodsCount%req.PageSize == 0 {
		extraPage = 0
	}
	res.Data = &dto.GoodsListRes{
		GoodsList:   ConvertToGoodsInfoList(goodsList),
		TotalPage:   goodsCount/req.PageSize + extraPage,
		TotalNumber: goodsCount,
		PageSize:    req.PageSize,
		Page:        req.Page,
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
		err := ss.storeService.AddGoods(ConvertFromGoodsInfo(req.Goods))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ss.storeService.UpdateGoods(ConvertFromGoodsInfo(req.Goods))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(storeServerLogTag, "RequestModifyStoreType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ss *StorehouseServer) RequestStoreTypeList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	list, err := ss.storeService.GetStoreTypeList()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.StoreTypeListRes{
		StoreTypeList: ConvertToStoreTypeInfoList(list),
	}
}

func (ss *StorehouseServer) RequestModifyStoreType(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyStoreTypeReq)
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ss.storeService.AddStoreType(ConvertFromStoreTypeInfo(req.StoreType))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ss.storeService.UpdateStoreType(ConvertFromStoreTypeInfo(req.StoreType))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(storeServerLogTag, "RequestModifyStoreType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}
