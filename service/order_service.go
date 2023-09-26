package service

import (
	"database/sql"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	orderServiceLogTag = "OrderService"

	extraPayAmount = 1.6
)

type ApplyPayOrderInfo struct {
	PayOrder  *model.PayOrderDao
	OrderList []*ApplyOrderInfo
}

type ApplyOrderInfo struct {
	Order *model.OrderDao
	Items []*model.OrderDetail
}

type OrderService struct {
	sqlCli             *sql.DB
	payOrderModel      *model.PayOrderModel
	orderModel         *model.OrderModel
	orderDetailModel   *model.OrderDetailModel
	orderDiscountModel *model.OrderDiscountModel
}

func NewOrderService(sqlCli *sql.DB) *OrderService {
	payOrderModel := model.NewPayOrderModel(sqlCli)
	orderModel := model.NewOrderModel(sqlCli)
	orderDetailModel := model.NewOrderDetailModel(sqlCli)
	orderDiscountModel := model.NewOrderDiscountModel(sqlCli)
	return &OrderService{
		payOrderModel:      payOrderModel,
		orderModel:         orderModel,
		orderDetailModel:   orderDetailModel,
		orderDiscountModel: orderDiscountModel,
		sqlCli:             sqlCli,
	}
}

func (os *OrderService) ApplyPayOrder(applyInfo *ApplyPayOrderInfo, dishMap map[uint32]*model.Dish,
	discountType uint8) (prepareID string, totalAmount, payAmount float64, err error) {
	tx, err := os.sqlCli.Begin()
	if err != nil {
		logger.Warn(orderServiceLogTag, "ApplyPayOrder Begin Failed|Err:%v", err)
		return
	}
	defer func() {
		if err != nil {
			rollErr := tx.Rollback()
			if rollErr != nil {
				logger.Warn(orderServiceLogTag, "ApplyPayOrder Rollback Failed|Err:%v", err)
				return
			}
		} else {
			err = tx.Commit()
			if err != nil {
				logger.Warn(orderServiceLogTag, "ApplyPayOrder Commit Failed|Err:%v", err)
				tx.Rollback()
			}
		}
	}()

	err = os.payOrderModel.InsertWithTx(tx, applyInfo.PayOrder)
	if err != nil {
		logger.Warn(orderServiceLogTag, "ApplyPayOrder InsertPayOrder Failed|Dao:%v|Err:%v", applyInfo.PayOrder, err)
		return
	}

	totalAmount, payAmount = float64(0), float64(0)
	discount, err := os.orderDiscountModel.GetDiscountByID(discountType)
	if err != nil {
		logger.Warn(orderServiceLogTag, "ApplyPayOrder GetDiscountByID|ID:%v|Err:%v", discountType, err)
	} else {
		discount = &model.OrderDiscount{}
	}

	for _, applyOrder := range applyInfo.OrderList {
		applyOrder.Order.PayOrderID = applyInfo.PayOrder.ID
		err = os.ApplyOrder(tx, applyOrder.Order, applyOrder.Items, dishMap, discount)
		if err != nil {
			logger.Warn(orderServiceLogTag, "ApplyPayOrder Failed|ID:%v|Err:%v", applyOrder.Order.ID, err)
			return
		}
		totalAmount += applyOrder.Order.TotalAmount
		payAmount += applyOrder.Order.PayAmount
	}
	applyInfo.PayOrder.TotalAmount = totalAmount
	applyInfo.PayOrder.PayAmount = payAmount
	err = os.payOrderModel.UpdatePayOrderInfoByID(tx, applyInfo.PayOrder, "total_amount", "pay_amount")
	if err != nil {
		logger.Warn(orderServiceLogTag, "UpdatePayOrderInfoByID Failed|ID:%v|Err:%v", applyInfo.PayOrder.ID, err)
		return
	}
	return
}

func (os *OrderService) ApplyOrder(tx *sql.Tx, order *model.OrderDao, items []*model.OrderDetail,
	dishMap map[uint32]*model.Dish, discount *model.OrderDiscount) error {
	totalAmount := float64(0)
	for _, item := range items {
		dish := dishMap[item.DishID]
		item.Price = dish.Price
		item.DishType = dish.DishType
		totalAmount = totalAmount + (item.Price * float64(item.Quantity))
	}
	payAmount := totalAmount
	payAmount = payAmount - discount.GetMealDiscount(order.MealType)
	if payAmount < 0 {
		payAmount = 0
	}

	order.TotalAmount = totalAmount
	order.PayAmount = payAmount

	err := os.orderModel.InsertWithTx(tx, order)
	if err != nil {
		logger.Warn(orderServiceLogTag, "Insert Order Failed|Err:%v", err)
		return err
	}

	for _, item := range items {
		item.OrderID = order.ID
	}
	err = os.orderDetailModel.BatchInsert(tx, items)
	if err != nil {
		logger.Warn(orderServiceLogTag, "BatchInsert OrderDetail Failed|Err:%v", err)
		return err
	}

	return nil
}

func (os *OrderService) GetPayOrderList(orderIDList []uint32, uid uint32, page, pageSize int32,
	orderStatus int8) ([]*model.PayOrderDao, int32, error) {
	orderList, err := os.payOrderModel.GetPayOrderList(orderIDList, uid, page, pageSize, orderStatus)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderList Failed|Err:%v", err)
		return nil, 0, err
	}

	if len(orderList) == 0 {
		return make([]*model.PayOrderDao, 0), 0, nil
	}

	orderCount, err := os.payOrderModel.GetPayOrderListCount(orderIDList, uid, page, pageSize, orderStatus)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetPayOrderListCount Failed|Err:%v", err)
		return nil, 0, err
	}

	return orderList, orderCount, nil
}

func (os *OrderService) GetOrderListByPayOrderID(payOrderList []uint32) ([]*model.OrderDao, map[uint32][]*model.OrderDetail, error) {
	orderList, err := os.orderModel.GetOrderListByPayOrder(payOrderList)
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

func (os *OrderService) GetOrderList(orderIDList []uint32, uid uint32, page, pageSize int32,
	orderStatus int8) ([]*model.OrderDao, int32, map[uint32][]*model.OrderDetail, error) {
	orderList, err := os.orderModel.GetOrderList(orderIDList, uid, page, pageSize, orderStatus)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}
	orderCount, err := os.orderModel.GetOrderListCount(orderIDList, uid, orderStatus)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderListCount Failed|Err:%v", err)
		return nil, 0, nil, err
	}

	if len(orderList) == 0 {
		return make([]*model.OrderDao, 0), 0, make(map[uint32][]*model.OrderDetail), nil
	}

	for _, order := range orderList {
		orderIDList = append(orderIDList, order.ID)
	}
	details, err := os.orderDetailModel.GetOrderDetailByOrderList(orderIDList)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderDetailByOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}

	detailMap := make(map[uint32][]*model.OrderDetail, 0)
	for _, detail := range details {
		if _, ok := detailMap[detail.OrderID]; ok == false {
			detailMap[detail.OrderID] = make([]*model.OrderDetail, 0)
		}
		detailMap[detail.OrderID] = append(detailMap[detail.OrderID], detail)
	}

	return orderList, orderCount, detailMap, nil
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

func (os *OrderService) LoginUserOrderDiscountInfo(uid uint32, discountType uint8) (float64, float64, error) {
	discountAmount := 0.0
	if discountType > 0 {
		discountInfo, err := os.orderDiscountModel.GetDiscountByID(discountType)
		if err != nil {
			logger.Warn(orderServiceLogTag, "GetDiscountByID Failed|Err:%v", err)
			return 0, 0, err
		}
		discountAmount = discountInfo.GetMealDiscount(enum.MealBreakfast)
	}

	orderList, err := os.orderModel.GetOrderList(make([]uint32, 0), uid, 1, 1, -1)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderList Failed|Err:%v", err)
		return 0, 0, err
	}
	minPay := extraPayAmount
	if len(orderList) > 0 {
		minPay = 0
	}

	return minPay, discountAmount, nil
}
