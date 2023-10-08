package service

import (
	"database/sql"
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
	storeTypeModel    *model.StorehouseTypeModel
	goodsModel        *model.GoodsModel
	goodsTypeModel    *model.GoodsTypeModel
	shoppingCartModel *model.ShoppingCartModel
	cartDetailModel   *model.CartDetailModel
}

func NewStoreService(sqlCli *sql.DB) *StoreService {
	storeTypeModel := model.NewStorehouseTypeModelWithDB(sqlCli)
	goodsModel := model.NewGoodsModelWithDB(sqlCli)
	goodsTypeModel := model.NewGoodsTypeModelWithDB(sqlCli)
	shoppingCartModel := model.NewShoppingCartModel(sqlCli)
	cartDetailModel := model.NewCartDetailModel(sqlCli)
	return &StoreService{
		storeTypeModel:    storeTypeModel,
		goodsModel:        goodsModel,
		goodsTypeModel:    goodsTypeModel,
		shoppingCartModel: shoppingCartModel,
		cartDetailModel:   cartDetailModel,
	}
}

func (ss *StoreService) Init() error {
	return nil
}

func (ss *StoreService) ReceivePurchase(details []*model.PurchaseDetail) {

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

func (ss *StoreService) GetGoodsTypeList() ([]*model.GoodsType, error) {
	typeList, err := ss.goodsTypeModel.GetGoodsTypes()
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetGoodsTypeList Failed|Err:%v", err)
		return nil, err
	}
	return typeList, nil
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
	err := ss.goodsTypeModel.UpdateGoodsType(goodsType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "UpdateGoodsType Failed|Err:%v", err)
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

	averagePrice = averagePrice * goodsType.Discount
	err = ss.goodsModel.UpdateGoodsPriceInfo(goodsID, averagePrice, priceMap)
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
