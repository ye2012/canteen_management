package server

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/service"
	"github.com/canteen_management/utils"

	"github.com/gin-gonic/gin"
)

type StorehouseServer struct {
	storeService *service.StoreService
}

func NewStorehouseServer(dbConf utils.Config) (*StorehouseServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(menuServerLogTag, "NewMenuServer Failed|Err:%v", err)
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
		logger.Warn(menuServerLogTag, "RequestModifyGoodsType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ss *StorehouseServer) RequestGoodsList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GoodsListReq)
	goodsList, err := ss.storeService.GoodsList(req.GoodsTypeID, req.StoreTypeID)
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.GoodsListRes{
		GoodsList: ConvertToGoodsInfoList(goodsList),
	}
}

func (ss *StorehouseServer) RequestModifyGoods(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyGoodsInfoReq)
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ss.storeService.AddGoods(ConvertFromGoodsInfo(req.Goods, ""))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ss.storeService.UpdateGoods(ConvertFromGoodsInfo(req.Goods, ""))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(menuServerLogTag, "RequestModifyStoreType Unknown OperateType|Type:%v", req.Operate)
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
		logger.Warn(menuServerLogTag, "RequestModifyStoreType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}
