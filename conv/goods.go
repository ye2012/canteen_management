package conv

import (
	"fmt"
	"github.com/canteen_management/dto"
	"github.com/canteen_management/model"
	"strconv"
	"strings"
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

func ConvertGoodsID(itemID string) (uint32, error) {
	ids := strings.Split(itemID, IndexDelimiter)
	if len(ids) != 2 {
		return 0, fmt.Errorf("id不合法")
	}

	goodsID, err := strconv.ParseInt(ids[1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("转化商品id失败|ID:%v", itemID)
	}
	return uint32(goodsID), nil
}

func ConvertGoodsListToGoodsNode(goodsMap map[uint32]*model.Goods, goodsTypes []*model.GoodsType,
	goodsSelectedMap map[string]float64) []*dto.GoodsNode {
	retData := make([]*dto.GoodsNode, 0)
	goodsTypeMap := make(map[uint32][]*model.Goods)
	for _, goods := range goodsMap {
		_, ok := goodsTypeMap[goods.GoodsTypeID]
		if ok == false {
			goodsTypeMap[goods.GoodsTypeID] = make([]*model.Goods, 0)
		}
		goodsTypeMap[goods.GoodsTypeID] = append(goodsTypeMap[goods.GoodsTypeID], goods)
	}
	for _, goodsType := range goodsTypes {
		goodsList, ok := goodsTypeMap[goodsType.ID]
		if !ok || len(goodsList) == 0 {
			continue
		}
		typeNode := &dto.GoodsNode{ID: fmt.Sprintf("%v", goodsType.ID), Name: goodsType.GoodsTypeName,
			Children: make([]*dto.GoodsNode, 0, len(goodsList))}
		typeSelected := int32(0)
		for _, goods := range goodsList {
			goodsNode := &dto.GoodsNode{
				ID:        fmt.Sprintf("%v_%v", goodsType.ID, goods.ID),
				Price:     goods.Price,
				Name:      goods.Name,
				GoodsID:   goods.ID,
				Left:      goods.Quantity,
				BatchSize: goods.BatchSize,
				BatchUnit: goods.BatchUnit,
			}
			typeNode.Children = append(typeNode.Children, goodsNode)
			if selected, ok := goodsSelectedMap[goodsNode.ID]; ok && selected > 0 {
				typeSelected += 1
			}
		}
		typeNode.SelectedNumber = typeSelected
		retData = append(retData, typeNode)
	}
	return retData
}
