package service

import (
	"database/sql"
	"fmt"

	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/utils"
)

const (
	inventoryServiceLogTag = "InventoryService"
)

type InventoryService struct {
	sqlCli               *sql.DB
	inventoryOrderModel  *model.InventoryOrderModel
	inventoryDetailModel *model.InventoryDetailModel
	goodsModel           *model.GoodsModel
	goodsHistoryModel    *model.GoodsHistoryModel
}

func NewInventoryService(sqlCli *sql.DB) *InventoryService {
	inventoryOrderModel := model.NewInventoryOrderModel(sqlCli)
	inventoryDetailModel := model.NewInventoryDetailModelWithDB(sqlCli)
	goodsModel := model.NewGoodsModelWithDB(sqlCli)
	goodsHistoryModel := model.NewGoodsHistoryModel(sqlCli)
	return &InventoryService{
		sqlCli:               sqlCli,
		inventoryOrderModel:  inventoryOrderModel,
		inventoryDetailModel: inventoryDetailModel,
		goodsModel:           goodsModel,
		goodsHistoryModel:    goodsHistoryModel,
	}
}

func (is *InventoryService) StartInventory(creator uint32) (*model.InventoryOrder, []*model.InventoryDetail, error) {
	goodsList, err := is.goodsModel.GetAllGoods()
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "StartInventory GetAllGoods Failed|Err:%v", err)
		return nil, nil, err
	}

	tx, err := is.sqlCli.Begin()
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "StartInventory Begin Failed|Err:%v", err)
		return nil, nil, err
	}
	defer utils.End(tx, err)

	inventory := &model.InventoryOrder{Creator: creator}
	err = is.inventoryOrderModel.InsertWithTx(tx, inventory)
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "StartInventory InsertInventoryOrder Failed|Err:%v", err)
		return nil, nil, err
	}

	detailList := make([]*model.InventoryDetail, 0, len(goodsList))
	for _, goods := range goodsList {
		detail := &model.InventoryDetail{
			InventoryID:  inventory.ID,
			GoodsID:      goods.ID,
			GoodsType:    goods.GoodsTypeID,
			ExpectNumber: goods.Quantity,
		}
		detailList = append(detailList, detail)
	}
	err = is.inventoryDetailModel.BatchInsert(detailList)
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "StartInventory InsertInventoryDetail Failed|Err:%v", err)
		return nil, nil, err
	}
	return inventory, detailList, nil
}

func (is *InventoryService) UpdateInventory(detail *model.InventoryDetail) error {
	err := is.inventoryDetailModel.UpdateDetailByCondition(detail, "id", "real_number", "tag", "status")
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "UpdateDetailByCondition Failed|Err:%v", err)
		return err
	}
	return nil
}

func (is *InventoryService) ApplyInventory(inventoryID uint32) error {
	details, err := is.inventoryDetailModel.GetDetail(inventoryID, enum.InventoryNew)
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "ApplyInventory GetDetail Failed|Err:%v", err)
		return err
	}
	if len(details) > 0 {
		return fmt.Errorf("还有未盘点的库存商品哦")
	}

	inventoryOrder := &model.InventoryOrder{ID: inventoryID, Status: enum.InventoryOrderFinish}
	err = is.inventoryOrderModel.UpdateInventoryOrderByCondition(inventoryOrder, "id", "status")
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "UpdateInventoryOrderByCondition Failed|Err:%v", err)
		return err
	}
	return nil
}

func (is *InventoryService) ReviewInventory(inventoryID uint32) error {
	tx, err := is.sqlCli.Begin()
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "StartInventory Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	inventoryOrder := &model.InventoryOrder{ID: inventoryID, Status: enum.InventoryOrderReviewed}
	err = is.inventoryOrderModel.UpdateInventoryOrderByConditionWithTx(tx, inventoryOrder, "id", "status")
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "UpdateInventoryOrderByCondition Failed|Err:%v", err)
		return err
	}

	details, err := is.inventoryDetailModel.GetDetail(inventoryID, enum.InventoryNeedFix)
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "ReviewInventory GetDetail Failed|Err:%v", err)
		return err
	}

	if len(details) > 0 {
		goodsIDList, goodsList := make([]uint32, 0, len(details)), make([]*model.Goods, 0, len(details))
		updateMap, historyList := make(map[uint32]float64), make([]*model.GoodsHistory, 0, len(details))
		for _, detail := range details {
			goodsIDList = append(goodsIDList, detail.GoodsID)
			updateMap[detail.GoodsID] = detail.RealNumber - detail.ExpectNumber
		}
		goodsList, err = is.goodsModel.GetGoodsByIDListWithLock(tx, goodsIDList)
		if err != nil {
			logger.Warn(inventoryServiceLogTag, "ReviewInventory GetGoodsByIDListWithLock Failed|Err:%v", err)
			return err
		}

		for _, goods := range goodsList {
			historyList = append(historyList,
				model.GenerateInventoryGoodsHistory(goods, updateMap[goods.ID], inventoryID))
			goods.Quantity = goods.Quantity + updateMap[goods.ID]
		}

		err = is.goodsModel.BatchUpdateQuantityWithTx(tx, goodsList)
		if err != nil {
			logger.Warn(inventoryServiceLogTag, "ReviewInventory BatchAddQuantity Failed|Err:%v", err)
			return err
		}
		err = is.goodsHistoryModel.BatchInsert(tx, historyList)
		if err != nil {
			logger.Warn(inventoryServiceLogTag, "ReviewInventory BatchInsertHistory Failed|Err:%v", err)
			return err
		}
	}
	return nil
}

func (is *InventoryService) InventoryOrderList(inventoryID, creator uint32, status int8, startTime, endTime int64,
	page, pageSize int32) ([]*model.InventoryOrder, int32, map[uint32][]*model.InventoryDetail, error) {
	inventoryList, err := is.inventoryOrderModel.GetInventoryOrderList(inventoryID, creator, status, startTime, endTime, page, pageSize)
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "InventoryOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}
	inventoryCount, err := is.inventoryOrderModel.GetInventoryOrderCount(inventoryID, creator, status, startTime, endTime)
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "GetInventoryOrderCount Failed|Err:%v", err)
		return nil, 0, nil, err
	}

	if len(inventoryList) == 0 {
		return make([]*model.InventoryOrder, 0), 0, make(map[uint32][]*model.InventoryDetail), nil
	}

	orderIDList := make([]uint32, 0, len(inventoryList))
	for _, outbound := range inventoryList {
		orderIDList = append(orderIDList, outbound.ID)
	}
	details, err := is.inventoryDetailModel.GetInventoryDetailByOrderList(orderIDList, 0)
	if err != nil {
		logger.Warn(inventoryServiceLogTag, "GetInventoryDetailByOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}

	detailMap := make(map[uint32][]*model.InventoryDetail, 0)
	for _, detail := range details {
		if _, ok := detailMap[detail.InventoryID]; ok == false {
			detailMap[detail.InventoryID] = make([]*model.InventoryDetail, 0)
		}
		detailMap[detail.InventoryID] = append(detailMap[detail.InventoryID], detail)
	}
	return inventoryList, inventoryCount, detailMap, nil
}
