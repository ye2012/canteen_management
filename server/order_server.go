package server

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/service"
	"github.com/canteen_management/utils"
	"github.com/gin-gonic/gin"
)

const (
	orderServerLogTag = "OrderServer"

	KeyPrice = "price"
)

type OrderServer struct {
	dishService  *service.DishService
	menuService  *service.MenuService
	orderService *service.OrderService
}

func NewOrderServer(dbConf utils.Config) (*OrderServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(orderServerLogTag, "NewOrderServer Failed|Err:%v", err)
		return nil, err
	}
	dishService := service.NewDishService(sqlCli)
	err = dishService.Init()
	if err != nil {
		return nil, err
	}
	menuService := service.NewMenuService(sqlCli)
	err = menuService.Init()
	if err != nil {
		return nil, err
	}
	orderService := service.NewOrderService(sqlCli)

	return &OrderServer{
		dishService:  dishService,
		menuService:  menuService,
		orderService: orderService,
	}, nil
}

func (os *OrderServer) RequestOrderMenu(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	typeMap, err := os.dishService.GetDishTypeMap()
	if err != nil {
		res.Code = enum.SystemError
		return
	}
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		res.Code = enum.SystemError
		return
	}

	nowTime := time.Now().Unix()
	orderDate := nowTime + 3600*24
	dayMenu, err := os.menuService.GetWeekMenuByTime(orderDate, 1)
	if err != nil {
		logger.Warn(orderServerLogTag, "RequestOrderMenu GetWeekMenuByTime Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}
	tomorrowData := ConvertMenuToOrderNode(orderDate, dayMenu, dishMap, typeMap)

	resData := dto.OrderMenuRes{}
	if nowTime < utils.GetMidDayTime(nowTime) {
		dayMenu, err = os.menuService.GetWeekMenuByTime(nowTime, 1)
		if err != nil {
			logger.Warn(orderServerLogTag, "RequestOrderMenu GetWeekMenuByTime Failed|Err:%v", err)
			res.Code = enum.SystemError
			return
		}

		delete(dayMenu, enum.MealBreakfast)
		delete(dayMenu, enum.MealLunch)

		resData = ConvertMenuToOrderNode(nowTime, dayMenu, dishMap, typeMap)
	}
	resData = append(resData, tomorrowData...)
	res.Data = resData
}

func ConvertMenuToOrderNode(menuDate int64, dayMenu map[uint8][]uint32, dishMap map[uint32]*model.Dish,
	typeMap map[uint32]*model.DishType) dto.OrderMenuRes {
	retData := dto.OrderMenuRes{}
	for mealType, totalDishList := range dayMenu {
		mealName := time.Unix(menuDate, 0).Format("01-02") + enum.GetMealName(mealType)
		retMeal := &dto.OrderNode{ID: fmt.Sprintf("%v_%v", menuDate, mealType), Name: mealName}
		dishListByType := make(map[uint32][]*model.Dish)
		for _, dishID := range totalDishList {
			dishType := dishMap[dishID].DishType
			if _, ok := dishListByType[dishType]; ok == false {
				dishListByType[dishType] = make([]*model.Dish, 0)
			}
			dishListByType[dishType] = append(dishListByType[dishType], dishMap[dishID])
		}

		retMeal.Children = make([]*dto.OrderNode, 0, len(dishListByType))
		for dishType, dishList := range dishListByType {
			retListByType := &dto.OrderNode{ID: fmt.Sprintf("%v", dishType), Name: typeMap[dishType].DishTypeName}
			retListByType.Children = make([]*dto.OrderNode, 0, len(dishList))
			for index, dish := range dishList {
				retDish := &dto.OrderNode{ID: fmt.Sprintf("%v_%v_%v", retMeal.ID, dish.ID, index),
					DishID: dish.ID, Name: dish.DishName, Price: dish.Price}
				retListByType.Children = append(retListByType.Children, retDish)
			}
			retMeal.Children = append(retMeal.Children, retListByType)
		}

		retData = append(retData, retMeal)
	}
	return retData
}

func (os *OrderServer) RequestApplyOrder(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ApplyPayOrderReq)
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	payOrder := &model.PayOrderDao{
		PrepareID:    "",
		Uid:          0,
		UnionID:      "UnionID",
		Address:      req.Address,
		DiscountType: 0,
		Status:       0,
	}

	applyPay := &service.ApplyPayOrderInfo{PayOrder: payOrder, OrderList: make([]*service.ApplyOrderInfo, 0)}
	for _, orderInfo := range req.OrderList {
		applyInfo := &service.ApplyOrderInfo{}
		orderDao := ConvertToOrderDao(0, orderInfo.ID, req.Address)
		if orderDao == nil {
			logger.Warn(orderServerLogTag, "Convert OrderDao Failed|Req:%#v", *req)
			continue
		}
		orderItems := ConvertToOrderDetailDao(orderInfo.OrderItems)

		applyInfo.Order = orderDao
		applyInfo.Items = orderItems
		applyPay.OrderList = append(applyPay.OrderList, applyInfo)
	}

	if len(applyPay.OrderList) == 0 {
		logger.Warn(orderServerLogTag, "ApplyList Length Zero")
		res.Code = enum.ParamsError
		return
	}

	prepareID, totalAmount, payAmount, err := os.orderService.ApplyPayOrder(applyPay, dishMap, 1)
	if err != nil {
		logger.Warn(orderServerLogTag, "ApplyPayOrder Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	req.TotalAmount = totalAmount
	req.PaymentAmount = payAmount

	resData := &dto.ApplyOrderRes{
		PayOrderInfo: req,
		PrepareID:    prepareID,
	}
	res.Data = resData
}

func ConvertToOrderDao(uid uint32, ID, addr string) *model.OrderDao {
	ids := strings.Split(ID, "_")
	if len(ids) != 2 {
		logger.Warn(orderServerLogTag, "ID illegal|ID:%v", ID)
		return nil
	}
	mealTime, _ := strconv.ParseInt(ids[0], 10, 32)
	mealType, _ := strconv.ParseInt(ids[1], 10, 32)

	return &model.OrderDao{
		OrderDate: time.Unix(mealTime, 0),
		MealType:  uint8(mealType),
		Uid:       uid,
		Address:   addr,
	}
}

func ConvertToOrderDetailDao(items []*dto.ApplyItem) []*model.OrderDetail {
	retList := make([]*model.OrderDetail, 0)
	for _, item := range items {
		retInfo := &model.OrderDetail{
			DishID:   item.DishID,
			Quantity: item.Quantity,
		}
		retList = append(retList, retInfo)
	}
	return retList
}

func (os *OrderServer) RequestPayOrderList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.PayOrderListReq)
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	payOrderList, payOrderNumber, err := os.orderService.GetPayOrderList([]uint32{}, req.Uid, req.Page, req.PageSize, req.OrderStatus)
	if err != nil {
		logger.Warn(orderServerLogTag, "GetPayOrderList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	payOrderIDList, payOrderIDMap := make([]uint32, 0), make(map[uint32]int)
	payOrderInfoList := make([]*dto.PayOrderInfo, 0)
	for _, payOrder := range payOrderList {
		payOrderIDList = append(payOrderIDList, payOrder.ID)
		payOrderInfo := &dto.PayOrderInfo{
			ID:            payOrder.ID,
			OrderList:     make([]*dto.OrderInfo, 0),
			Address:       payOrder.Address,
			TotalAmount:   payOrder.TotalAmount,
			PaymentAmount: payOrder.PayAmount,
			Status:        payOrder.Status,
		}
		payOrderIDMap[payOrder.ID] = len(payOrderInfoList)
		payOrderInfoList = append(payOrderInfoList, payOrderInfo)
	}

	orderList, detailMap, err := os.orderService.GetOrderListByPayOrderID(payOrderIDList)
	if err != nil {
		logger.Warn(orderServerLogTag, "GetOrderListByPayOrderID Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	orderInfoList := ConvertToOrderInfoList(orderList, detailMap, dishMap)
	for _, orderInfo := range orderInfoList {
		if index, ok := payOrderIDMap[orderInfo.PayOrderID]; ok {
			payOrderInfoList[index].OrderList = append(payOrderInfoList[index].OrderList, orderInfo)
		}
	}
	resData := &dto.PayOrderListRes{
		OrderList: payOrderInfoList,
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: payOrderNumber,
		},
	}
	res.Data = resData
}

func (os *OrderServer) RequestOrderList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.OrderListReq)
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}
	orderIDList := make([]uint32, 0)
	if req.OrderID > 0 {
		orderIDList = append(orderIDList, req.OrderID)
	}
	orderList, totalNumber, detailMap, err := os.orderService.GetOrderList(orderIDList, req.Uid, req.Page, req.PageSize, req.OrderStatus)
	if err != nil {
		logger.Warn(orderServerLogTag, "GetOrderList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	orderInfoList := ConvertToOrderInfoList(orderList, detailMap, dishMap)
	resData := &dto.OrderListRes{OrderList: orderInfoList,
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: totalNumber,
		},
	}
	res.Data = resData
}

func ConvertToOrderInfoList(orderList []*model.OrderDao, detailMap map[uint32][]*model.OrderDetail,
	dishMap map[uint32]*model.Dish) []*dto.OrderInfo {
	retList := make([]*dto.OrderInfo, 0)
	for _, order := range orderList {
		retInfo := &dto.OrderInfo{
			ID:            fmt.Sprintf("%v_%v", order.OrderDate.Unix(), order.MealType),
			Name:          order.OrderDate.Format("01-02") + enum.GetMealName(order.MealType),
			OrderID:       fmt.Sprintf("%v", order.ID),
			OrderNo:       "",
			PayOrderID:    order.PayOrderID,
			Address:       order.Address,
			TotalAmount:   order.TotalAmount,
			PaymentAmount: order.PayAmount,
			OrderItems:    make([]*dto.ApplyItem, 0),
			CreateTime:    order.CreateAt.Unix(),
			OrderStatus:   order.Status,
		}
		if details, ok := detailMap[order.ID]; ok {
			orderItems := make([]*dto.ApplyItem, 0, len(details))
			for _, detail := range details {
				item := &dto.ApplyItem{
					DishID:   detail.DishID,
					DishName: dishMap[detail.DishID].DishName,
					Price:    detail.Price,
					Quantity: detail.Quantity,
				}
				orderItems = append(orderItems, item)
			}
			retInfo.OrderItems = orderItems
		}
		retList = append(retList, retInfo)
	}
	return retList
}

func (os *OrderServer) RequestDiscountList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	discountList, err := os.orderService.GetOrderDiscountList()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetOrderDiscountList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	discountInfoList := make([]*dto.OrderDiscountInfo, 0, len(discountList))
	for _, discount := range discountList {
		err = discount.ConvertDiscount()
		if err != nil {
			logger.Warn(orderServerLogTag, "ConvertDiscount Failed|Err:%v", err)
			continue
		}

		discountInfo := &dto.OrderDiscountInfo{
			ID:                discount.ID,
			DiscountTypeName:  discount.DiscountTypeName,
			BreakfastDiscount: discount.GetMealDiscount(enum.MealBreakfast),
			LunchDiscount:     discount.GetMealDiscount(enum.MealLunch),
			DinnerDiscount:    discount.GetMealDiscount(enum.MealDinner),
		}
		discountInfoList = append(discountInfoList, discountInfo)
	}
	retData := &dto.OrderDiscountListRes{
		DiscountList: discountInfoList,
	}
	res.Data = retData
}

func (os *OrderServer) RequestModifyDiscount(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyOrderDiscountReq)

	discountMap := make(map[uint8]float64)
	discountMap[enum.MealBreakfast] = req.DiscountInfo.BreakfastDiscount
	discountMap[enum.MealLunch] = req.DiscountInfo.LunchDiscount
	discountMap[enum.MealDinner] = req.DiscountInfo.DinnerDiscount
	discount := &model.OrderDiscount{
		ID:               req.DiscountInfo.ID,
		DiscountTypeName: req.DiscountInfo.DiscountTypeName,
	}
	err := discount.FromDiscountMap(discountMap)
	if err != nil {
		logger.Warn(orderServerLogTag, "FromDiscountMap Failed|Err:%v", err)
		res.Code = enum.ParamsError
		return
	}

	switch req.Operate {
	case enum.OperateTypeAdd:
		err = os.orderService.AddOrderDiscount(discount)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err = os.orderService.UpdateOrderDiscount(discount)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(orderServerLogTag, "RequestModifyDiscount Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}
