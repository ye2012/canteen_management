package conv

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/model"
)

func ConvertToInventoryInfoList(inventoryList []*model.InventoryOrder, detailMap map[uint32][]*model.InventoryDetail,
	goodsMap map[uint32]*model.Goods) []*dto.InventoryOrderInfo {
	retList := make([]*dto.InventoryOrderInfo, 0, len(inventoryList))
	for _, inventory := range inventoryList {
		retInfo := &dto.InventoryOrderInfo{
			ID:        inventory.ID,
			GoodsList: make([]*dto.InventoryGoodsInfo, 0),
			Status:    inventory.Status,
		}
		details, ok := detailMap[inventory.ID]
		if !ok {
			continue
		}
		for _, detail := range details {
			inventoryGoods := &dto.InventoryGoodsInfo{
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
			retInfo.GoodsList = append(retInfo.GoodsList, inventoryGoods)
		}
		retList = append(retList, retInfo)
	}
	return retList
}

func ConvertFromApplyInventory(inventory *dto.InventoryGoodsInfo) *model.InventoryDetail {
	return &model.InventoryDetail{
		ID:         inventory.ID,
		RealNumber: inventory.RealNumber,
		Tag:        inventory.Tag,
		Status:     inventory.Status,
	}

}
