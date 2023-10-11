package conv

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/model"
)

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
			ID:            purchaseGoods.ID,
			GoodsID:       goods.ID,
			GoodsType:     goods.GoodsTypeID,
			ExpectNumber:  purchaseGoods.ExpectNumber,
			ReceiveNumber: purchaseGoods.ReceiveNumber,
			Price:         goods.Price,
		}
		detailList = append(detailList, detail)
	}
	return detailList
}

func ConvertToPurchaseInfoList(purchaseList []*model.PurchaseOrder, detailMap map[uint32][]*model.PurchaseDetail,
	goodsMap map[uint32]*model.Goods, supplierMap map[uint32]*model.Supplier) []*dto.PurchaseOrderInfo {
	retList := make([]*dto.PurchaseOrderInfo, 0, len(purchaseList))
	for _, purchase := range purchaseList {
		retInfo := &dto.PurchaseOrderInfo{
			ID:            purchase.ID,
			Supplier:      purchase.Supplier,
			SupplierName:  supplierMap[purchase.Supplier].Name,
			GoodsList:     make([]*dto.PurchaseGoodsInfo, 0),
			TotalAmount:   purchase.TotalAmount,
			PaymentAmount: purchase.PayAmount,
			Status:        uint8(purchase.Status),
		}
		details, ok := detailMap[purchase.ID]
		if !ok {
			continue
		}
		for _, detail := range details {
			purchaseGoods := &dto.PurchaseGoodsInfo{
				PurchaseGoodsBase: dto.PurchaseGoodsBase{
					ID:           detail.ID,
					GoodsID:      detail.GoodsID,
					Name:         goodsMap[detail.GoodsID].Name,
					GoodsTypeID:  detail.GoodsType,
					ExpectNumber: detail.ExpectNumber,
				},
				ReceiveNumber: detail.ReceiveNumber,
				Price:         detail.Price,
			}
			retInfo.GoodsList = append(retInfo.GoodsList, purchaseGoods)
		}
		retList = append(retList, retInfo)
	}
	return retList
}

func ConvertFromApplyOutbound(goodsList []*dto.OutboundGoodsInfo, goodsMap map[uint32]*model.Goods) []*model.OutboundDetail {
	detailList := make([]*model.OutboundDetail, 0, len(goodsList))
	for _, outboundGoods := range goodsList {
		goods, ok := goodsMap[outboundGoods.GoodsID]
		if ok == false {
			continue
		}
		detail := &model.OutboundDetail{
			ID:        outboundGoods.ID,
			GoodsID:   goods.ID,
			GoodsType: goods.GoodsTypeID,
			OutNumber: outboundGoods.ExpectNumber,
			Price:     goods.Price,
		}
		detailList = append(detailList, detail)
	}
	return detailList
}

func ConvertToOutboundInfoList(outboundList []*model.OutboundOrder, detailMap map[uint32][]*model.OutboundDetail,
	goodsMap map[uint32]*model.Goods) []*dto.OutboundOrderInfo {
	retList := make([]*dto.OutboundOrderInfo, 0, len(outboundList))
	for _, outbound := range outboundList {
		retInfo := &dto.OutboundOrderInfo{
			ID:          outbound.ID,
			GoodsList:   make([]*dto.OutboundGoodsInfo, 0),
			TotalAmount: outbound.TotalAmount,
		}
		details, ok := detailMap[outbound.ID]
		if !ok {
			continue
		}
		for _, detail := range details {
			purchaseGoods := &dto.OutboundGoodsInfo{
				PurchaseGoodsBase: dto.PurchaseGoodsBase{
					ID:           detail.ID,
					GoodsID:      detail.GoodsID,
					Name:         goodsMap[detail.GoodsID].Name,
					GoodsTypeID:  detail.GoodsType,
					ExpectNumber: detail.OutNumber,
				},
			}
			retInfo.GoodsList = append(retInfo.GoodsList, purchaseGoods)
		}
		retList = append(retList, retInfo)
	}
	return retList
}