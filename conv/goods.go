package conv

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/model"
)

func ConvertFromStoreTypeInfo(info *dto.StoreTypeInfo) *model.StorehouseType {
	return &model.StorehouseType{ID: info.StoreTypeID, StoreTypeName: info.StoreTypeName}
}

func ConvertToStoreTypeInfoList(daoList []*model.StorehouseType) []*dto.StoreTypeInfo {
	retList := make([]*dto.StoreTypeInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.StoreTypeInfo{StoreTypeID: dao.ID, StoreTypeName: dao.StoreTypeName})
	}
	return retList
}

func ConvertFromGoodsTypeInfo(info *dto.GoodsTypeInfo) *model.GoodsType {
	return &model.GoodsType{ID: info.GoodsTypeID, GoodsTypeName: info.GoodsTypeName, Discount: info.Discount}
}

func ConvertToGoodsTypeInfoList(daoList []*model.GoodsType) []*dto.GoodsTypeInfo {
	retList := make([]*dto.GoodsTypeInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.GoodsTypeInfo{GoodsTypeID: dao.ID, GoodsTypeName: dao.GoodsTypeName, Discount: dao.Discount})
	}
	return retList
}

func ConvertFromGoodsInfo(info *dto.GoodsInfo) *model.Goods {
	return &model.Goods{ID: info.GoodsID, Name: info.GoodsName, GoodsTypeID: info.GoodsType, StoreTypeID: info.StoreType,
		Picture: info.Picture, BatchSize: info.BatchSize, BatchUnit: info.BatchUnit, Price: info.Price, Quantity: info.Quantity}
}

func ConvertToGoodsInfoList(daoList []*model.Goods) []*dto.GoodsInfo {
	retList := make([]*dto.GoodsInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.GoodsInfo{GoodsID: dao.ID, GoodsName: dao.Name, GoodsType: dao.GoodsTypeID,
			StoreType: dao.StoreTypeID, Picture: dao.Picture, BatchSize: dao.BatchSize, BatchUnit: dao.BatchUnit,
			Price: dao.Price, Quantity: dao.Quantity})
	}
	return retList
}

func ConvertToGoodsPriceList(daoList []*model.Goods) []*dto.GoodsPriceInfo {
	retList := make([]*dto.GoodsPriceInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.GoodsPriceInfo{GoodsID: dao.ID, GoodsName: dao.Name,
			PriceList: dao.ToGoodsPrice(), AveragePrice: dao.Price})
	}
	return retList
}
