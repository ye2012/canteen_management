package conv

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/model"
)

func ConvertToInventoryInfoList(inventoryList []*model.InventoryOrder, detailMap map[uint32][]*model.InventoryDetail,
	goodsMap map[uint32]*model.Goods, goodsTypes []*model.GoodsType) []*dto.InventoryOrderInfo {
	retList := make([]*dto.InventoryOrderInfo, 0, len(inventoryList))
	for _, inventory := range inventoryList {
		retInfo := &dto.InventoryOrderInfo{
			ID:        inventory.ID,
			GoodsList: make([]*dto.InventoryGoodsNode, 0),
			Status:    inventory.Status,
		}
		details, ok := detailMap[inventory.ID]
		if !ok {
			continue
		}
		totalCount, exceptionCount, exceptionAmount := int32(0), int32(0), float64(0)
		for _, detail := range details {
			inventoryGoods := &dto.InventoryGoodsNode{
				PurchaseGoodsBase: dto.PurchaseGoodsBase{
					ID:           detail.ID,
					GoodsID:      detail.GoodsID,
					Name:         goodsMap[detail.GoodsID].Name,
					GoodsTypeID:  detail.GoodsType,
					ExpectNumber: detail.ExpectNumber,
				},
				RealNumber: detail.RealNumber,
				Tag:        detail.Tag,
				Status:     detail.Status,
			}
			totalCount++
			if detail.Status == enum.InventoryNeedFix {
				exceptionCount++
				exceptionAmount += (detail.RealNumber - detail.ExpectNumber) * goodsMap[detail.GoodsID].Price
			}
			retInfo.GoodsList = append(retInfo.GoodsList, inventoryGoods)
		}
		retInfo.TotalCount = totalCount
		retInfo.ExceptionCount = exceptionCount
		retInfo.ExceptionAmount = exceptionAmount
		retList = append(retList, retInfo)
	}
	return retList
}

func ConvertFromApplyInventory(inventory *dto.InventoryGoodsNode) *model.InventoryDetail {
	return &model.InventoryDetail{
		ID:         inventory.ID,
		RealNumber: inventory.RealNumber,
		Tag:        inventory.Tag,
		Status:     inventory.Status,
	}

}

func ConvertGoodsListToInventoryNode(inventory *model.InventoryOrder, details []*model.InventoryDetail,
	goodsMap map[uint32]*model.Goods, goodsTypes []*model.GoodsType) *dto.InventoryOrderInfo {
	retInfo := &dto.InventoryOrderInfo{
		ID:        inventory.ID,
		GoodsList: make([]*dto.InventoryGoodsNode, 0),
		Status:    inventory.Status,
	}
	goodsTypeMap := make(map[uint32][]*model.InventoryDetail)
	for _, goods := range details {
		_, ok := goodsTypeMap[goods.GoodsType]
		if ok == false {
			goodsTypeMap[goods.GoodsType] = make([]*model.InventoryDetail, 0)
		}
		goodsTypeMap[goods.GoodsType] = append(goodsTypeMap[goods.GoodsType], goods)
	}
	for _, goodsType := range goodsTypes {
		goodsList, ok := goodsTypeMap[goodsType.ID]
		if !ok || len(goodsList) == 0 {
			continue
		}
		typeNode := &dto.InventoryGoodsNode{PurchaseGoodsBase: dto.PurchaseGoodsBase{ID: goodsType.ID, Name: goodsType.GoodsTypeName},
			Children: make([]*dto.InventoryGoodsNode, 0, len(goodsList))}
		for _, detail := range goodsList {
			inventoryGoods := &dto.InventoryGoodsNode{
				PurchaseGoodsBase: dto.PurchaseGoodsBase{
					ID:           detail.ID,
					GoodsID:      detail.GoodsID,
					Name:         goodsMap[detail.GoodsID].Name,
					GoodsTypeID:  detail.GoodsType,
					ExpectNumber: detail.ExpectNumber,
				},
				BatchSize:  goodsMap[detail.GoodsID].BatchSize,
				BatchUnit:  goodsMap[detail.GoodsID].BatchUnit,
				RealNumber: detail.RealNumber,
				Tag:        detail.Tag,
				Status:     detail.Status,
			}
			typeNode.Children = append(typeNode.Children, inventoryGoods)
		}
		retInfo.GoodsList = append(retInfo.GoodsList, typeNode)
	}
	return retInfo
}
