package service

import (
	"database/sql"
	"fmt"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/utils"
	"math"
	"time"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	storeServiceLogTag = "StoreService"
)

type StoreService struct {
	sqlCli              *sql.DB
	storeTypeModel      *model.StorehouseTypeModel
	goodsModel          *model.GoodsModel
	goodsTypeModel      *model.GoodsTypeModel
	goodsHistoryModel   *model.GoodsHistoryModel
	shoppingCartModel   *model.ShoppingCartModel
	cartDetailModel     *model.CartDetailModel
	outboundModel       *model.OutboundOrderModel
	outboundDetailModel *model.OutboundDetailModel
}

func NewStoreService(sqlCli *sql.DB) *StoreService {
	storeTypeModel := model.NewStorehouseTypeModelWithDB(sqlCli)
	goodsModel := model.NewGoodsModelWithDB(sqlCli)
	goodsTypeModel := model.NewGoodsTypeModelWithDB(sqlCli)
	goodsHistoryModel := model.NewGoodsHistoryModel(sqlCli)
	shoppingCartModel := model.NewShoppingCartModel(sqlCli)
	cartDetailModel := model.NewCartDetailModel(sqlCli)
	outboundModel := model.NewOutboundOrderModelWithDB(sqlCli)
	outboundDetailModel := model.NewOutboundDetailModelWithDB(sqlCli)
	return &StoreService{
		sqlCli:              sqlCli,
		storeTypeModel:      storeTypeModel,
		goodsModel:          goodsModel,
		goodsTypeModel:      goodsTypeModel,
		goodsHistoryModel:   goodsHistoryModel,
		shoppingCartModel:   shoppingCartModel,
		cartDetailModel:     cartDetailModel,
		outboundModel:       outboundModel,
		outboundDetailModel: outboundDetailModel,
	}
}

func (ss *StoreService) Init() error {
	return nil
}

func (ss *StoreService) GetStoreTypeList() ([]*model.StorehouseType, error) {
	typeList, err := ss.storeTypeModel.GetStorehouseTypes()
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetStoreTypeList Failed|Err:%v", err)
		return nil, err
	}
	return typeList, nil
}

func (ss *StoreService) AddStoreType(storeType *model.StorehouseType) error {
	err := ss.storeTypeModel.Insert(storeType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "Insert StoreType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) UpdateStoreType(storeType *model.StorehouseType) error {
	err := ss.storeTypeModel.UpdateStorehouseType(storeType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "UpdateStoreType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) DeleteStoreType(storeTypeID uint32) error {
	count, err := ss.goodsModel.GetGoodsCount(0, storeTypeID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "DeleteStoreType GetGoodsCount Failed|Err:%v", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("该仓库下还有商品，无法删除")
	}
	err = ss.storeTypeModel.DeleteStorehouseType(storeTypeID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "UpdateStoreType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) GetGoodsTypeList() ([]*model.GoodsType, error) {
	typeList, err := ss.goodsTypeModel.GetGoodsTypes()
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsTypeList Failed|Err:%v", err)
		return nil, err
	}
	return typeList, nil
}

func (ss *StoreService) GetGoodsTypeMap() (map[uint32]*model.GoodsType, error) {
	retMap := make(map[uint32]*model.GoodsType)
	typeList, err := ss.goodsTypeModel.GetGoodsTypes()
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsTypeList Failed|Err:%v", err)
		return nil, err
	}
	for _, typeInfo := range typeList {
		retMap[typeInfo.ID] = typeInfo
	}
	return retMap, nil
}

func (ss *StoreService) AddGoodsType(goodsType *model.GoodsType) error {
	err := ss.goodsTypeModel.Insert(goodsType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "Insert GoodsType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) UpdateGoodsType(goodsType *model.GoodsType) error {
	tx, err := ss.sqlCli.Begin()
	if err != nil {
		logger.Warn(storeServiceLogTag, "UpdateGoodsType Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	preType, err := ss.goodsTypeModel.GetGoodsTypesByID(goodsType.ID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsTypesByID Failed|Err:%v", err)
		return err
	}

	if preType.Discount != goodsType.Discount {
		goodsList, err := ss.goodsModel.GetGoods(goodsType.ID, 0, 1, 10000)
		if err != nil {
			logger.Warn(storeServiceLogTag, "UpdateGoodsType GetGoods Failed|Err:%v", err)
			return err
		}
		for _, goods := range goodsList {
			goods.Price = goods.AveragePrice * goodsType.Discount
		}
		ss.goodsModel.BatchUpdateByTagWithTx(tx, goodsList, "price")
	}

	err = ss.goodsTypeModel.UpdateGoodsTypeWithTx(tx, goodsType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "UpdateGoodsType Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) DeleteGoodsType(goodsTypeID uint32) error {
	count, err := ss.goodsModel.GetGoodsCount(goodsTypeID, 0)
	if err != nil {
		logger.Warn(storeServiceLogTag, "DeleteGoodsType GetGoodsCount Failed|Err:%v", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("该商品类型下还有商品，无法删除")
	}
	err = ss.goodsTypeModel.DeleteGoodsType(goodsTypeID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsTypesByID Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) GetGoodsMap() (map[uint32]*model.Goods, error) {
	goodsList, err := ss.goodsModel.GetAllGoods()
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetAllGoods Failed|Err:%v", err)
		return nil, err
	}

	goodsMap := make(map[uint32]*model.Goods)
	for _, goods := range goodsList {
		goodsMap[goods.ID] = goods
	}
	return goodsMap, nil
}

func (ss *StoreService) GoodsList(goodsType, storeType uint32, page, pageSize int32) ([]*model.Goods, int32, error) {
	goodsList, err := ss.goodsModel.GetGoods(goodsType, storeType, page, pageSize)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoods Failed|Err:%v", err)
		return nil, 0, err
	}

	goodsCount, err := ss.goodsModel.GetGoodsCount(goodsType, storeType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsCount Failed|Err:%v", err)
		return nil, 0, err
	}

	return goodsList, goodsCount, nil
}

func (ss *StoreService) AddGoods(goods *model.Goods) error {
	err := ss.goodsModel.Insert(goods)
	if err != nil {
		logger.Warn(storeServiceLogTag, "Insert Goods Failed|Err:%v", err)
		return err
	}
	if goods.Quantity > 0 {
		err = ss.goodsHistoryModel.BatchInsert(nil, []*model.GoodsHistory{model.GenerateInitGoodsHistory(goods)})
		logger.Warn(storeServiceLogTag, "Insert GoodsHistory Failed|Err:%v", err)
	}
	return nil
}

func (ss *StoreService) UpdateGoods(goods *model.Goods) error {
	err := ss.goodsModel.UpdateGoodsInfo(goods)
	if err != nil {
		logger.Warn(storeServiceLogTag, "UpdateGoods Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) DeleteGoods(goodsID uint32) error {
	err := ss.goodsModel.DeleteGoods(goodsID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "DeleteGoods Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) UpdateGoodsPrice(goodsID uint32, priceMap map[uint8]float64) error {
	averagePrice, count := 0.0, 0
	for _, price := range priceMap {
		if math.Abs(price) < 0.0000001 {
			continue
		}
		averagePrice += price
		count += 1
	}
	if count <= 0 {
		return nil
	}
	averagePrice /= float64(count)

	goods, err := ss.goodsModel.GetGoodsByID(goodsID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsByID Failed|Err:%v", err)
		return err
	}

	goodsType, err := ss.goodsTypeModel.GetGoodsTypesByID(goods.GoodsTypeID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsTypesByID Failed|Err:%v", err)
		return err
	}

	finalPrice := averagePrice * goodsType.Discount
	err = ss.goodsModel.UpdateGoodsPriceInfo(goodsID, averagePrice, finalPrice, priceMap)
	if err != nil {
		logger.Warn(storeServiceLogTag, "UpdateGoodsPrice Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) GetCart(uid uint32, cartType enum.CartType) (*model.ShoppingCart, []*model.CartDetail, error) {
	carts, err := ss.shoppingCartModel.GetCart(cartType, uid)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetCart Failed|Err:%v", err)
		return nil, nil, err
	}

	cart, cartDetails := (*model.ShoppingCart)(nil), make([]*model.CartDetail, 0)
	if len(carts) > 0 {
		if carts[0].CreateAt.Unix() < utils.GetZeroTime(time.Now().Unix()) {
			err = ss.shoppingCartModel.Delete(cartType, uid, 0)
			if err != nil {
				logger.Warn(storeServiceLogTag, "Delete ShoppingCart Failed|Err:%v", err)
				return nil, nil, err
			}
			cartIDs := make([]uint32, 0, len(carts))
			for _, preCart := range carts {
				cartIDs = append(cartIDs, preCart.ID)
			}
			err = ss.cartDetailModel.Delete(cartIDs)
			if err != nil {
				logger.Warn(storeServiceLogTag, "Delete CartDetail Failed|Err:%v", err)
				return nil, nil, err
			}
		} else {
			cart = carts[0]
			cartDetails, err = ss.cartDetailModel.GetCartDetail(cart.ID)
			if err != nil {
				logger.Warn(storeServiceLogTag, "GetCartDetail Failed|Err:%v", err)
				return nil, nil, err
			}
		}
	}
	return cart, cartDetails, nil
}

func (ss *StoreService) ApplyOutboundOrder(outbound *model.OutboundOrder, details []*model.OutboundDetail) (err error) {
	tx, err := ss.sqlCli.Begin()
	if err != nil {
		logger.Warn(storeServiceLogTag, "ApplyOutboundOrder Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	totalAmount := 0.0
	for _, item := range details {
		totalAmount += item.Price * item.OutNumber
	}
	outbound.TotalAmount = totalAmount
	err = ss.outboundModel.InsertWithTx(tx, outbound)
	if err != nil {
		logger.Warn(storeServiceLogTag, "Insert Outbound Failed|Err:%v", err)
		return err
	}
	for _, item := range details {
		item.OutboundID = outbound.ID
	}
	err = ss.outboundDetailModel.BatchInsertWithTx(tx, details)
	if err != nil {
		logger.Warn(storeServiceLogTag, "BatchInsert OutboundDetail Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) ReviewOutboundOrder(outboundID uint32) error {
	tx, err := ss.sqlCli.Begin()
	if err != nil {
		logger.Warn(storeServiceLogTag, "FinishOutboundOrder Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	order, err := ss.outboundModel.GetOutboundOrderWithLock(tx, outboundID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetOutboundOrderWithLock Failed|Err:%v", err)
		return err
	}
	if order == nil {
		logger.Warn(storeServiceLogTag, "ReviewOutboundOrder Not Found|ID:%v|Err:%v", outboundID, err)
		return fmt.Errorf("订单未找到")
	}
	if order.Status != enum.OutboundNew {
		return fmt.Errorf("订单状态异常")
	}

	order.Status = enum.OutboundReviewed
	err = ss.outboundModel.UpdateOutboundWithTx(tx, order, "status")
	if err != nil {
		logger.Warn(storeServiceLogTag, "Review UpdateOutboundWithTx Failed|Err:%v", err)
		return err
	}
	return err
}

func (ss *StoreService) FinishOutboundOrder(outboundID uint32) error {
	tx, err := ss.sqlCli.Begin()
	if err != nil {
		logger.Warn(storeServiceLogTag, "FinishOutboundOrder Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	order, err := ss.outboundModel.GetOutboundOrderWithLock(tx, outboundID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetOutboundOrderWithLock Failed|Err:%v", err)
		return err
	}
	if order == nil {
		logger.Warn(storeServiceLogTag, "ReviewOutboundOrder Not Found|ID:%v|Err:%v", outboundID, err)
		return fmt.Errorf("订单未找到")
	}
	if order.Status != enum.OutboundReviewed {
		return fmt.Errorf("订单状态异常")
	}

	details, err := ss.outboundDetailModel.GetDetail(outboundID)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetOutboundOrder Detail Failed|Err:%v", err)
		return err
	}

	err = ss.finishOutboundOrder(tx, order, details)
	if err == nil {
		order.Status = enum.OutboundFinish
		err = ss.outboundModel.UpdateOutboundWithTx(tx, order, "status")
		if err != nil {
			logger.Warn(storeServiceLogTag, "UpdateOutboundStatus Failed|Err:%v", err)
			return err
		}
	}
	return err
}

func (ss *StoreService) finishOutboundOrder(tx *sql.Tx, outbound *model.OutboundOrder, details []*model.OutboundDetail) (err error) {
	goodsIDList, outMap := make([]uint32, 0, len(details)), make(map[uint32]float64)
	for _, item := range details {
		goodsIDList = append(goodsIDList, item.GoodsID)
		outMap[item.GoodsID] = item.OutNumber
	}
	goodsList, err := ss.goodsModel.GetGoodsByIDListWithLock(tx, goodsIDList)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsByIDListWithLock Failed|Err:%v", err)
		return err
	}
	updateMap, historyList := make(map[uint32]float64), make([]*model.GoodsHistory, 0, len(details))
	for _, goods := range goodsList {
		outNumber, ok := outMap[goods.ID]
		if !ok {
			continue
		}
		if outNumber > goods.Quantity {
			logger.Warn(storeServiceLogTag, "OutNumber Extent StoreQuantity|Goods:%v|Out:%v|Store:%v",
				goods.ID, outNumber, goods.Quantity)
			return fmt.Errorf("%v出库数量超出库存数量", goods.Name)
		}
		updateMap[goods.ID] = outNumber
	}
	for _, item := range details {
		item.OutboundID = outbound.ID
	}

	for _, goods := range goodsList {
		historyList = append(historyList,
			model.GenerateOutboundGoodsHistory(goods, -updateMap[goods.ID], outbound.ID))
		goods.Quantity = goods.Quantity - updateMap[goods.ID]
	}

	err = ss.goodsModel.BatchUpdateQuantityWithTx(tx, goodsList)
	if err != nil {
		logger.Warn(storeServiceLogTag, "ApplyOutboundOrder BatchAddQuantity Failed|Err:%v", err)
		return err
	}
	err = ss.goodsHistoryModel.BatchInsert(tx, historyList)
	if err != nil {
		logger.Warn(storeServiceLogTag, "ApplyOutboundOrder BatchInsertHistory Failed|Err:%v", err)
		return err
	}
	return nil
}

func (ss *StoreService) GetOutboundList(uid, outboundID uint32, startTime, endTime int64, status int8,
	page, pageSize int32) ([]*model.OutboundOrder, int32, map[uint32][]*model.OutboundDetail, error) {
	outboundList, err := ss.outboundModel.GetOutboundOrderList(outboundID, uid, startTime, endTime, status, page, pageSize)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetOutboundOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}
	outboundCount, err := ss.outboundModel.GetOutboundOrderCount(outboundID, uid, startTime, endTime, status)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetOutboundOrderCount Failed|Err:%v", err)
		return nil, 0, nil, err
	}

	if len(outboundList) == 0 {
		return make([]*model.OutboundOrder, 0), 0, make(map[uint32][]*model.OutboundDetail), nil
	}

	orderIDList := make([]uint32, 0, len(outboundList))
	for _, outbound := range outboundList {
		orderIDList = append(orderIDList, outbound.ID)
	}
	details, err := ss.outboundDetailModel.GetOutboundDetailByOrderList(orderIDList, 0)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetOutboundDetailByOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}

	detailMap := make(map[uint32][]*model.OutboundDetail, 0)
	for _, detail := range details {
		if _, ok := detailMap[detail.OutboundID]; ok == false {
			detailMap[detail.OutboundID] = make([]*model.OutboundDetail, 0)
		}
		detailMap[detail.OutboundID] = append(detailMap[detail.OutboundID], detail)
	}
	return outboundList, outboundCount, detailMap, nil
}

func (ss *StoreService) GetGoodsHistoryList(goodsID, changeType uint32, startTime, endTime int64,
	page, pageSize int32) ([]*model.GoodsHistory, int32, error) {
	history, err := ss.goodsHistoryModel.GetGoodsHistory(goodsID, changeType, startTime, endTime, page, pageSize)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsHistoryList Failed|Err:%v", err)
		return nil, 0, err
	}
	count, err := ss.goodsHistoryModel.GetGoodsHistoryCount(goodsID, changeType, startTime, endTime)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsHistoryCount Failed|Err:%v", err)
		return nil, 0, err
	}
	return history, count, nil
}
