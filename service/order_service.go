package service

import (
	"database/sql"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	orderServiceLogTag = "OrderService"
)

type OrderService struct {
	sqlCli             *sql.DB
	orderModel         *model.OrderModel
	orderDetailModel   *model.OrderDetailModel
	orderDiscountModel *model.OrderDiscountModel
	orderUserModel     *model.OrderUserModel
}

func NewOrderService(sqlCli *sql.DB) *OrderService {
	orderModel := model.NewOrderModel(sqlCli)
	orderDetailModel := model.NewOrderDetailModel(sqlCli)
	orderDiscountModel := model.NewOrderDiscountModel(sqlCli)
	orderUserModel := model.NewOrderUserModel(sqlCli)
	return &OrderService{
		orderModel:         orderModel,
		orderDetailModel:   orderDetailModel,
		orderDiscountModel: orderDiscountModel,
		orderUserModel:     orderUserModel,
		sqlCli:             sqlCli,
	}
}

func (os *OrderService) ApplyOrder(order *model.OrderDao, items []*model.OrderDetail,
	dishMap map[uint32]*model.Dish, discountType uint8) (prepareID string, err error) {
	tx, err := os.sqlCli.Begin()
	if err != nil {
		logger.Warn(orderServiceLogTag, "ApplyOrder Begin Failed|Err:%v", err)
		return
	}
	defer func() {
		if err != nil {
			rollErr := tx.Rollback()
			if rollErr != nil {
				logger.Warn(orderServiceLogTag, "ApplyOrder Rollback Failed|Err:%v", err)
				return
			}
		} else {
			err = tx.Commit()
			if err != nil {
				logger.Warn(orderServiceLogTag, "ApplyOrder Commit Failed|Err:%v", err)
				tx.Rollback()
			}
		}
	}()

	err = os.orderModel.InsertWithTx(tx, order)
	if err != nil {
		logger.Warn(orderServiceLogTag, "Insert Order Failed|Err:%v", err)
		return
	}

	totalAmount := float64(0)
	for _, item := range items {
		dish := dishMap[item.DishID]
		item.Price = dish.Price
		item.DishType = dish.DishType
		item.OrderID = order.ID
		totalAmount = totalAmount + (item.Price * float64(item.Quantity))
	}

	err = os.orderDetailModel.BatchInsert(tx, items)
	if err != nil {
		logger.Warn(orderServiceLogTag, "BatchInsert OrderDetail Failed|Err:%v", err)
		return
	}

	payAmount := totalAmount
	discount, err := os.orderDiscountModel.GetDiscountByID(discountType)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetDiscountByID|ID:%v|Err:%v", discountType, err)
	} else {
		payAmount -= discount.GetMealDiscount(order.MealType)
		if payAmount < 0 {
			payAmount = 0
		}
	}

	order.TotalAmount = totalAmount
	order.PayAmount = payAmount
	// todo prepare id

	err = os.orderModel.UpdateOrderInfoByID(tx, order, "total_amount", "pay_amount")
	if err != nil {
		logger.Warn(orderServiceLogTag, "UpdateOrderInfoByID Failed|Err:%v", err)
		return
	}
	return
}

func (os *OrderService) GetOrderList(orderID, uid, page, pageSize uint32,
	orderStatus int8) ([]*model.OrderDao, map[uint32][]*model.OrderDetail, error) {
	orderList, err := os.orderModel.GetOrderList(orderID, uid, page, pageSize, orderStatus)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderList Failed|Err:%v", err)
		return nil, nil, err
	}

	if len(orderList) == 0 {
		return make([]*model.OrderDao, 0), make(map[uint32][]*model.OrderDetail), nil
	}

	orderIDList := make([]uint32, 0)
	for _, order := range orderList {
		orderIDList = append(orderIDList, order.ID)
	}
	details, err := os.orderDetailModel.GetOrderDetailByOrderList(orderIDList)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderDetailByOrderList Failed|Err:%v", err)
		return nil, nil, err
	}

	detailMap := make(map[uint32][]*model.OrderDetail, 0)
	for _, detail := range details {
		if _, ok := detailMap[detail.OrderID]; ok == false {
			detailMap[detail.OrderID] = make([]*model.OrderDetail, 0)
		}
		detailMap[detail.OrderID] = append(detailMap[detail.OrderID], detail)
	}

	return orderList, detailMap, nil
}

func (os *OrderService) GetOrder(orderID uint32) (*model.OrderDao, []*model.OrderDetail, error) {
	orderInfo, err := os.orderModel.GetOrder(orderID)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrder Failed|Err:%v", err)
		return nil, nil, err
	}

	details, err := os.orderDetailModel.GetOrderDetail(orderID)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderDetail Failed|Err:%v", err)
		return nil, nil, err
	}

	return orderInfo, details, nil
}

func (os *OrderService) GetOrderDetail(orderID uint32) ([]*model.OrderDetail, error) {
	detail, err := os.orderDetailModel.GetOrderDetail(orderID)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderDetail Failed|Err:%v", err)
		return nil, err
	}
	return detail, nil
}

func (os *OrderService) GetOrderUserList(phoneNumber string, discountLevel int32, page, pageSize uint32) ([]*model.OrderUser, uint32, error) {
	userList, err := os.orderUserModel.GetOrderUser(phoneNumber, discountLevel, page, pageSize)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderUser Failed|Err:%v", err)
		return nil, 0, err
	}

	userNumber, err := os.orderUserModel.GetOrderUserCount(phoneNumber, discountLevel)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderUserCount Failed|Err:%v", err)
		return nil, 0, err
	}
	return userList, uint32(userNumber), nil
}

func (os *OrderService) AddOrderUser(userList []*model.OrderUser) error {
	err := os.orderUserModel.BatchInsert(userList)
	if err != nil {
		logger.Warn(orderServiceLogTag, "AddOrderUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (os *OrderService) UpdateOrderUser(userInfo *model.OrderUser) error {
	err := os.orderUserModel.UpdateOrderUser(userInfo)
	if err != nil {
		logger.Warn(orderServiceLogTag, "ModifyOrderUser Failed|Err:%v", err)
		return err
	}
	return nil
}

func (os *OrderService) GetOrderDiscountList() ([]*model.OrderDiscount, error) {
	discountList, err := os.orderDiscountModel.GetDiscountList()
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderDiscountList Failed|Err:%v", err)
		return nil, err
	}
	return discountList, nil
}

func (os *OrderService) AddOrderDiscount(discountInfo *model.OrderDiscount) error {
	return os.orderDiscountModel.Insert(discountInfo)
}

func (os *OrderService) UpdateOrderDiscount(discountInfo *model.OrderDiscount) error {
	return os.orderDiscountModel.UpdateDiscountType(discountInfo)
}

//func GenerateOrder(id uint32) string {
//	now := time.Now().Format("200602011504")
//	return fmt.Sprintf("O%03d%v", id%1000, now)
//}
