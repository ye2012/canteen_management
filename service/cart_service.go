package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/utils"
)

const (
	cartServiceLogTag = "CartService"
)

type CartService struct {
	sqlCli            *sql.DB
	shoppingCartModel *model.ShoppingCartModel
	cartDetailModel   *model.CartDetailModel
}

func NewCartService(sqlCli *sql.DB) *CartService {
	shoppingCartModel := model.NewShoppingCartModel(sqlCli)
	cartDetailModel := model.NewCartDetailModel(sqlCli)
	return &CartService{
		sqlCli:            sqlCli,
		shoppingCartModel: shoppingCartModel,
		cartDetailModel:   cartDetailModel,
	}
}

func (cs *CartService) CheckCart(cartID, uid uint32, cartType enum.CartType) error {
	cart, err := cs.shoppingCartModel.GetCartByID(cartID)
	if err != nil {
		logger.Warn(cartServiceLogTag, "GetCartByID Failed|Err:%v", err)
		return err
	}
	if cart.CartType != cartType {
		logger.Warn(cartServiceLogTag, "CartType Not Match|Cart:%v|ReqType:%v", cart.CartType, cartType)
		return fmt.Errorf("购物车类型不匹配")
	}
	if cart.Uid != uid {
		logger.Warn(cartServiceLogTag, "Uid Not Match|Uid:%v|ReqUid:%v", cart.Uid, uid)
		return fmt.Errorf("购物车用户不匹配")
	}
	return nil
}

func (cs *CartService) ClearCart(uid uint32, cartType enum.CartType) error {
	tx, err := cs.sqlCli.Begin()
	if err != nil {
		logger.Warn(cartServiceLogTag, "ClearCart Begin Failed|Err:%v", err)
		return err
	}
	defer utils.End(tx, err)

	cartList, err := cs.shoppingCartModel.GetCartByTxWithLock(tx, cartType, uid)
	if err != nil {
		logger.Warn(cartServiceLogTag, "GetCartByTxWithLock Failed|Uid:%v|Err:%v", uid, err)
		return err
	}

	err = cs.shoppingCartModel.DeleteByTx(tx, cartType, uid, 0)
	if err != nil {
		logger.Warn(cartServiceLogTag, "Delete Cart Failed|Uid:%v|Err:%v", uid, err)
		return err
	}

	cartIDList := make([]uint32, 0)
	for _, cart := range cartList {
		cartIDList = append(cartIDList, cart.ID)
	}
	if len(cartIDList) == 0 {
		logger.Info(cartServiceLogTag, "Empty Cart|Uid:%v", uid)
		return nil
	}

	err = cs.cartDetailModel.DeleteWithTx(tx, cartIDList)
	if err != nil {
		logger.Warn(cartServiceLogTag, "Delete CartDetail Failed|IDs:%#v|Err:%v", cartIDList, err)
		return err
	}
	return nil
}

func (cs *CartService) GetCart(uid uint32, cartType enum.CartType) (*model.ShoppingCart, []*model.CartDetail, error) {
	carts, err := cs.shoppingCartModel.GetCart(cartType, uid)
	if err != nil {
		logger.Warn(cartServiceLogTag, "GetCart Failed|Err:%v", err)
		return nil, nil, err
	}

	cart, cartDetails := (*model.ShoppingCart)(nil), make([]*model.CartDetail, 0)
	if len(carts) > 0 {
		if carts[0].CreateAt.Unix() < utils.GetZeroTime(time.Now().Unix()) {
			err = cs.shoppingCartModel.Delete(cartType, uid, 0)
			if err != nil {
				logger.Warn(cartServiceLogTag, "Delete ShoppingCart Failed|Err:%v", err)
				return nil, nil, err
			}
			cartIDs := make([]uint32, 0, len(carts))
			for _, preCart := range carts {
				cartIDs = append(cartIDs, preCart.ID)
			}
			err = cs.cartDetailModel.Delete(cartIDs)
			if err != nil {
				logger.Warn(cartServiceLogTag, "Delete CartDetail Failed|Err:%v", err)
				return nil, nil, err
			}
		} else {
			cart = carts[0]
			cartDetails, err = cs.cartDetailModel.GetCartDetail(cart.ID)
			if err != nil {
				logger.Warn(cartServiceLogTag, "GetCartDetail Failed|Err:%v", err)
				return nil, nil, err
			}
		}
	}
	return cart, cartDetails, nil
}

func (cs *CartService) ModifyCart(uid uint32, itemID string, quantity float64,
	cartType enum.CartType) (*model.ShoppingCart, []*model.CartDetail, error) {
	tx, err := utils.Begin(cs.sqlCli)
	if err != nil {
		logger.Warn(cartServiceLogTag, "ModifyCart Begin Failed|Err:%v", err)
		return nil, nil, err
	}
	defer utils.End(tx, err)

	carts, err := cs.shoppingCartModel.GetCartByTxWithLock(tx, cartType, uid)
	if err != nil {
		logger.Warn(cartServiceLogTag, "GetCard Failed|Err:%v", err)
		return nil, nil, err
	}

	cart, cartDetails := (*model.ShoppingCart)(nil), make([]*model.CartDetail, 0)
	if len(carts) > 0 {
		if carts[0].CreateAt.Unix() < utils.GetZeroTime(time.Now().Unix()) {
			err = cs.shoppingCartModel.DeleteWithTx(tx, cartType, uid, 0)
			if err != nil {
				logger.Warn(cartServiceLogTag, "Delete ShoppingCart Failed|Err:%v", err)
				return nil, nil, err
			}
			cartIDs := make([]uint32, 0, len(carts))
			for _, preCart := range carts {
				cartIDs = append(cartIDs, preCart.ID)
			}
			err = cs.cartDetailModel.DeleteWithTx(tx, cartIDs)
			if err != nil {
				logger.Warn(cartServiceLogTag, "Delete CartDetail Failed|Err:%v", err)
				return nil, nil, err
			}
		} else {
			cart = carts[0]
			cartDetails, err = cs.cartDetailModel.GetCartDetailWithLock(tx, cart.ID)
			if err != nil {
				logger.Warn(cartServiceLogTag, "GetCartDetailWithLock Failed|Err:%v", err)
				return nil, nil, err
			}
		}
	}
	if cart == nil {
		cart = &model.ShoppingCart{
			CartType: cartType,
			Uid:      uid,
		}
		err = cs.shoppingCartModel.InsertWithTx(tx, cart)
		if err != nil {
			logger.Warn(cartServiceLogTag, "Insert ShoppingCart Failed|Err:%v", err)
			return nil, nil, err
		}
	}

	found := false
	for _, detail := range cartDetails {
		if detail.ItemID == itemID {
			detail.Quantity = quantity
			err = cs.cartDetailModel.UpdateDetail(tx, detail)
			if err != nil {
				logger.Warn(cartServiceLogTag, "UpdateDetail Failed|Err:%v", err)
				return nil, nil, err
			}
			found = true
		}
	}
	if found == false {
		detail := &model.CartDetail{CartID: cart.ID, ItemID: itemID, Quantity: quantity}
		err = cs.cartDetailModel.BatchInsert(tx, []*model.CartDetail{detail})
		if err != nil {
			logger.Warn(cartServiceLogTag, "Insert CartDetail Failed|Err:%v", err)
			return nil, nil, err
		}
		cartDetails = append(cartDetails, detail)
	}

	return cart, cartDetails, nil
}
