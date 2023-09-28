package service

import (
	"database/sql"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/utils"
	"time"
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
	discountInfo, err := os.orderDiscountModel.GetDiscountByID(discountType)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetDiscountByID Failed|ID:%v|Err:%v", discountType, err)
		return
	}
	discountAmount := discountInfo.GetMealDiscount(enum.MealBreakfast)

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

	timeStart, timeEnd := utils.GetDayTimeRange(applyInfo.OrderList[0].Order.OrderDate.Unix())
	prePayOrders, err := os.payOrderModel.GetPayOrderListWithLock(tx, []uint32{}, applyInfo.PayOrder.Uid,
		[]int8{enum.PayOrderNew, enum.PayOrderFinish}, timeStart, timeEnd)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetPayOrderListWithLock Failed|Dao:%v|Err:%v", applyInfo.PayOrder, err)
		return
	}
	extraPay := extraPayAmount
	for _, prePay := range prePayOrders {
		extraPay = 0
		discountAmount -= prePay.DiscountAmount
	}
	if discountAmount < 0 {
		discountAmount = 0
	}

	err = os.payOrderModel.InsertWithTx(tx, applyInfo.PayOrder)
	if err != nil {
		logger.Warn(orderServiceLogTag, "ApplyPayOrder InsertPayOrder Failed|Dao:%v|Err:%v", applyInfo.PayOrder, err)
		return
	}

	realDiscount := float64(0)
	for _, applyOrder := range applyInfo.OrderList {
		applyOrder.Order.PayOrderID = applyInfo.PayOrder.ID
		err = os.ApplyOrder(tx, applyOrder.Order, applyOrder.Items, dishMap, discountAmount, extraPay)
		if err != nil {
			logger.Warn(orderServiceLogTag, "ApplyPayOrder Failed|ID:%v|Err:%v", applyOrder.Order.ID, err)
			return
		}
		totalAmount += applyOrder.Order.TotalAmount
		payAmount += applyOrder.Order.PayAmount
		realDiscount += applyOrder.Order.DiscountAmount
		discountAmount -= applyOrder.Order.DiscountAmount
		extraPay = 0
	}
	applyInfo.PayOrder.TotalAmount = totalAmount
	applyInfo.PayOrder.PayAmount = payAmount
	applyInfo.PayOrder.DiscountAmount = realDiscount
	err = os.payOrderModel.UpdatePayOrderInfoByID(tx, applyInfo.PayOrder, "total_amount",
		"pay_amount", "discount_amount")
	if err != nil {
		logger.Warn(orderServiceLogTag, "UpdatePayOrderInfoByID Failed|ID:%v|Err:%v", applyInfo.PayOrder.ID, err)
		return
	}
	return
}

func (os *OrderService) CancelPayOrder(orderID uint32) (err error) {
	order := &model.PayOrderDao{ID: orderID, Status: enum.PayOrderCancel}
	err = os.payOrderModel.UpdatePayOrderInfoByID(nil, order, "status")
	if err != nil {
		logger.Warn(orderServiceLogTag, "CancelPayOrder Failed|Dao:%v|Err:%v", order, err)
		return
	}
	return
}

func (os *OrderService) DeliverOrder(orderID uint32) (err error) {
	order := &model.OrderDao{ID: orderID, Status: enum.OrderFinish, DeliverTime: time.Now()}
	err = os.orderModel.UpdateOrderInfoByID(nil, order, "status", "deliver_time")
	if err != nil {
		logger.Warn(orderServiceLogTag, "UpdateOrderInfoByID Failed|Dao:%v|Err:%v", order, err)
		return
	}
	return
}

func (os *OrderService) ApplyOrder(tx *sql.Tx, order *model.OrderDao, items []*model.OrderDetail,
	dishMap map[uint32]*model.Dish, discountAmount, extraPay float64) error {
	totalAmount := float64(0)
	for _, item := range items {
		dish := dishMap[item.DishID]
		item.Price = dish.Price
		item.DishType = dish.DishType
		totalAmount = totalAmount + (item.Price * float64(item.Quantity))
	}
	payAmount := totalAmount
	payAmount = payAmount - discountAmount
	realDiscount := discountAmount
	if payAmount < 0 {
		realDiscount = payAmount
		payAmount = 0
	}

	order.TotalAmount = totalAmount
	order.PayAmount = payAmount + extraPay
	order.DiscountAmount = realDiscount

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

	orderCount, err := os.payOrderModel.GetPayOrderListCount(orderIDList, uid, orderStatus)
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

func (os *OrderService) GetFloors(buildingID uint32, status int8, startTime, endTime int64, mealType uint8) ([]int32, error) {
	floors, err := os.orderModel.GetFloors(buildingID, status, startTime, endTime, mealType)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetFloors Failed|Err:%v", err)
		return floors, err
	}
	return floors, err
}

func (os *OrderService) GetOrderList(orderIDList []uint32, uid uint32, buildingID, floor uint32, room string,
	orderStatus int8, page, pageSize int32, startTime, endTime int64) ([]*model.OrderDao, int32, map[uint32][]*model.OrderDetail, error) {
	orderList, err := os.orderModel.GetOrderList(orderIDList, uid, buildingID, floor, room, orderStatus,
		startTime, endTime, page, pageSize)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}
	orderCount, err := os.orderModel.GetOrderListCount(orderIDList, uid, orderStatus, buildingID, floor, room, startTime, endTime)
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

	timeStart, timeEnd := utils.GetDayTimeRange(time.Now().Add(time.Hour * 24).Unix())
	payOrderList, err := os.payOrderModel.GetAllPayOrderList(make([]uint32, 0), uid,
		[]int8{enum.PayOrderNew, enum.PayOrderFinish}, timeStart, timeEnd)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderList Failed|Err:%v", err)
		return 0, 0, err
	}
	minPay := extraPayAmount
	for _, payOrder := range payOrderList {
		discountAmount -= payOrder.DiscountAmount
		if payOrder.Status == enum.PayOrderFinish {
			minPay = 0
		}
	}

	return minPay, discountAmount, nil
}
