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
	orderUserModel     *model.OrderUserModel
}

func NewOrderService(sqlCli *sql.DB) *OrderService {
	payOrderModel := model.NewPayOrderModel(sqlCli)
	orderModel := model.NewOrderModel(sqlCli)
	orderDetailModel := model.NewOrderDetailModel(sqlCli)
	orderDiscountModel := model.NewOrderDiscountModel(sqlCli)
	orderUserModel := model.NewOrderUserModel(sqlCli)
	return &OrderService{
		payOrderModel:      payOrderModel,
		orderModel:         orderModel,
		orderDetailModel:   orderDetailModel,
		orderDiscountModel: orderDiscountModel,
		orderUserModel:     orderUserModel,
		sqlCli:             sqlCli,
	}
}

func (os *OrderService) ApplyPayOrder(applyInfo *ApplyPayOrderInfo, dishMap map[uint32]*model.Dish,
	discountType uint8, cartID uint32) (prepareID string, totalAmount, payAmount float64, err error) {
	discountAmount := 0.0
	if discountType > 0 {
		discountInfo := &model.OrderDiscount{}
		discountInfo, err = os.orderDiscountModel.GetDiscountByID(discountType)
		if err != nil {
			logger.Warn(orderServiceLogTag, "GetDiscountByID Failed|ID:%v|Err:%v", discountType, err)
			return
		}
		discountAmount = discountInfo.GetMealDiscount(enum.MealBreakfast)
	}

	tx, err := os.sqlCli.Begin()
	if err != nil {
		logger.Warn(orderServiceLogTag, "ApplyPayOrder Begin Failed|Err:%v", err)
		return
	}
	defer utils.End(tx, err)

	if cartID > 0 {
		prePayOrder := (*model.PayOrderDao)(nil)
		prePayOrder, err = os.payOrderModel.GetPayOrderByCondition(tx, " WHERE `cart_id` = ?", cartID)
		if err != nil && err != sql.ErrNoRows {
			logger.Warn(orderServiceLogTag, "ApplyPayOrder GetPayOrderByCartID Failed|Err:%v", err)
			return
		}
		if prePayOrder != nil {
			logger.Warn(orderServiceLogTag, "ApplyPayOrder Cart Already Processed|CartID:%v", cartID)
			return "", 0, 0, fmt.Errorf("订单已经提交了")
		}
	}

	timeStart, timeEnd := utils.GetDayTimeRange(applyInfo.OrderList[0].Order.OrderDate.Unix())
	prePayOrders, err := os.payOrderModel.GetPayOrderListWithLock(tx, []uint32{}, applyInfo.PayOrder.Uid,
		[]int8{enum.PayOrderNew, enum.PayOrderFinish}, timeStart, timeEnd)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetPayOrderListWithLock Failed|Dao:%v|Err:%v", applyInfo.PayOrder, err)
		return
	}
	extraPay := extraPayAmount
	for _, prePay := range prePayOrders {
		if prePay.Status == enum.PayOrderFinish {
			extraPay = 0
		}
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

func (os *OrderService) CancelPayOrder(orderID uint32, payMethod uint8) (err error) {
	payOrder, err := os.payOrderModel.GetPayOrder(orderID)
	if err != nil {
		logger.Warn(orderServiceLogTag, "CancelPayOrder Get Failed|ID:%v|Err:%v", orderID, err)
		return err
	}
	if payOrder.PayMethod != payMethod {
		logger.Warn(orderServiceLogTag, "CancelPayOrder PayMethod Not Match|ID:%v|Err:%v", orderID, err)
		return fmt.Errorf("订单类型不匹配")
	}

	payOrder.Status = enum.PayOrderCancel
	err = os.payOrderModel.UpdatePayOrderInfoByID(nil, payOrder, "status")
	if err != nil {
		logger.Warn(orderServiceLogTag, "CancelPayOrder Failed|Dao:%v|Err:%v", payOrder, err)
		return
	}

	order := &model.OrderDao{PayOrderID: orderID, Status: enum.OrderCancel}
	os.orderModel.UpdateOrderInfo(nil, order, "pay_order_id", "status")
	return
}

func (os *OrderService) FinishPayOrder(orderID uint32, payMethod uint8) (err error) {
	payOrder, err := os.payOrderModel.GetPayOrder(orderID)
	if err != nil {
		logger.Warn(orderServiceLogTag, "FinishPayOrder Get Failed|ID:%v|Err:%v", orderID, err)
		return err
	}
	if payOrder.PayMethod != payMethod {
		logger.Warn(orderServiceLogTag, "FinishPayOrder PayMethod Not Match|ID:%v|Err:%v", orderID, err)
		return fmt.Errorf("订单类型不匹配")
	}

	payOrder.Status = enum.PayOrderFinish
	err = os.payOrderModel.UpdatePayOrderInfoByID(nil, payOrder, "status")
	if err != nil {
		logger.Warn(orderServiceLogTag, "CancelPayOrder Failed|Dao:%v|Err:%v", payOrder, err)
		return
	}

	order := &model.OrderDao{PayOrderID: orderID, Status: enum.OrderPaid}
	os.orderModel.UpdateOrderInfo(nil, order, "pay_order_id", "status")
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
		totalAmount += item.Price * float64(item.Quantity)
	}
	payAmount := totalAmount - discountAmount
	realDiscount := discountAmount
	if payAmount < 0 {
		realDiscount = totalAmount
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
	details, err := os.orderDetailModel.GetOrderDetailByOrderList(orderIDList, 0, 0)
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
	floors, err := os.orderModel.GetFloors(buildingID, status, -1, startTime, endTime, mealType)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetFloors Failed|Err:%v", err)
		return floors, err
	}
	return floors, err
}

func (os *OrderService) GetAllOrder(mealType uint8, startTime, endTime int64, status int8,
	dishType, dishID uint32) ([]*model.OrderDao, map[uint32][]*model.OrderDetail, error) {
	orderList, err := os.orderModel.GetAllOrder(mealType, startTime, endTime, status, -1)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetAllOrder Failed|Err:%v", err)
		return nil, nil, err
	}

	if len(orderList) == 0 {
		return make([]*model.OrderDao, 0), make(map[uint32][]*model.OrderDetail), nil
	}

	orderIDList := make([]uint32, 0)
	for _, order := range orderList {
		orderIDList = append(orderIDList, order.ID)
	}
	details, err := os.orderDetailModel.GetOrderDetailByOrderList(orderIDList, dishType, dishID)
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

func (os *OrderService) GetOrderList(orderIDList []uint32, uid uint32, mealType uint8, buildingID, floor uint32, room string,
	orderStatus, payMethod int8, page, pageSize int32, startTime, endTime int64) ([]*model.OrderDao, int32, map[uint32][]*model.OrderDetail, error) {
	orderList, err := os.orderModel.GetOrderList(orderIDList, uid, mealType, buildingID, floor, room, orderStatus, payMethod,
		startTime, endTime, page, pageSize)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderList Failed|Err:%v", err)
		return nil, 0, nil, err
	}
	orderCount, err := os.orderModel.GetOrderListCount(orderIDList, uid, mealType, orderStatus, payMethod, buildingID, floor, room, startTime, endTime)
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
	details, err := os.orderDetailModel.GetOrderDetailByOrderList(orderIDList, 0, 0)
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

func (os *OrderService) DeleteOrderDiscount(id uint32) error {
	condition := " WHERE `discount_level` = ? "
	userList, err := os.orderUserModel.GetOrderUserByCondition(condition, id)
	if err != nil {
		logger.Warn(orderServiceLogTag, "DeleteOrderDiscount GetOrderUser Failed|Err:%v", err)
		return err
	}
	if len(userList) > 0 {
		return fmt.Errorf("该折扣类型下还有人员，无法删除")
	}

	err = os.orderDiscountModel.DeleteDiscountType(id)
	if err != nil {
		logger.Warn(orderServiceLogTag, "DeleteOrderDiscount Failed|Err:%v", err)
		return err
	}

	return nil
}

func (os *OrderService) LoginUserOrderDiscountInfo(uid uint32, discountType uint8) (float64, float64, float64, error) {
	discountAmount, totalDiscount := 0.0, 0.0
	if discountType > 0 {
		discountInfo, err := os.orderDiscountModel.GetDiscountByID(discountType)
		if err != nil {
			logger.Warn(orderServiceLogTag, "GetDiscountByID Failed|Err:%v", err)
			return 0, 0, 0, err
		}
		discountAmount = discountInfo.GetMealDiscount(enum.MealBreakfast)
		totalDiscount = discountAmount
	}

	timeStart, timeEnd := utils.GetDayTimeRange(time.Now().Add(time.Hour * 24).Unix())
	payOrderList, err := os.payOrderModel.GetAllPayOrderList(make([]uint32, 0), uid,
		[]int8{enum.PayOrderNew, enum.PayOrderFinish}, timeStart, timeEnd)
	if err != nil {
		logger.Warn(orderServiceLogTag, "GetOrderList Failed|Err:%v", err)
		return 0, 0, 0, err
	}
	minPay := extraPayAmount
	for _, payOrder := range payOrderList {
		discountAmount -= payOrder.DiscountAmount
		if payOrder.Status == enum.PayOrderFinish {
			minPay = 0
		}
		logger.Debug(orderServiceLogTag, "Discount:%f|OrderDiscount:%f", discountAmount, payOrder.DiscountAmount)
	}

	return minPay, totalDiscount, discountAmount, nil
}
