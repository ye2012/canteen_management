package server

import (
	"strings"
	"time"

	"github.com/canteen_management/conv"
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
)

type OrderServer struct {
	dishService  *service.DishService
	menuService  *service.MenuService
	orderService *service.OrderService
	userService  *service.UserService
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
	userService := service.NewUserService(sqlCli)

	return &OrderServer{
		dishService:  dishService,
		menuService:  menuService,
		orderService: orderService,
		userService:  userService,
	}, nil
}

func (os *OrderServer) RequestOrderMenu(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.OrderMenuReq)
	uid := req.Uid
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
	orderDate := utils.GetZeroTime(nowTime + 3600*24)
	dayMenu, err := os.menuService.GetWeekMenuByTime(orderDate, 1)
	if err != nil {
		logger.Warn(orderServerLogTag, "RequestOrderMenu GetWeekMenuByTime Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}
	dishQuantityMap, totalCost, totalGoods := make(map[string]float64), 0.0, 0.0
	if uid != 0 {
		_, cartDetails, err := os.orderService.GetCart(uid)
		if err != nil {
			res.Code = enum.SystemError
			return
		}
		for _, detail := range cartDetails {
			dishQuantityMap[detail.ItemID] = detail.Quantity
			dishID, _ := conv.ConvertDishID(detail.ItemID)
			dish, ok := dishMap[dishID]
			if ok {
				totalCost += dish.Price * detail.Quantity
			}
			totalGoods += detail.Quantity
		}
	}
	tomorrowData := conv.ConvertMenuToOrderNode(orderDate, dayMenu, dishMap, typeMap, dishQuantityMap, true)

	resData := dto.OrderMenuRes{
		Menu:       tomorrowData,
		GoodsMap:   dishQuantityMap,
		TotalGoods: totalGoods,
		TotalCost:  totalCost,
	}
	res.Data = resData
}

func (os *OrderServer) RequestApplyOrder(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ApplyPayOrderReq)
	uid := uint32(1)

	prepareID, code, msg := os.ProcessApplyOrder(uid, req, enum.PayMethodWeChat)
	if code != enum.Success {
		res.Code = code
		res.Msg = msg
		return
	}

	resData := &dto.ApplyOrderRes{
		PayOrderInfo: req,
		PrepareID:    prepareID,
	}
	res.Data = resData
}

func (os *OrderServer) ProcessApplyOrder(uid uint32, req *dto.ApplyPayOrderReq, payMethod enum.PayMethod) (string, enum.ErrorCode, string) {
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetDishIDMap Failed|Err:%v", err)
		return "", enum.SystemError, ""
	}

	wxUser, err := os.userService.GetWxUser(uid)
	if err != nil || wxUser == nil {
		logger.Warn(orderServerLogTag, "GetWxUser Failed|Err:%v", err)
		return "", enum.SystemError, "用户不存在"
	}
	discountLevel := os.userService.GetWxUserDiscount(wxUser.OpenID)

	payOrder := &model.PayOrderDao{
		PrepareID:      "",
		Uid:            uid,
		OpenID:         wxUser.OpenID,
		BuildingID:     req.BuildingID,
		Floor:          req.Floor,
		Room:           req.Room,
		DiscountAmount: 0,
		PayMethod:      payMethod,
		Status:         enum.PayOrderNew,
		PayAmount:      req.PaymentAmount,
	}

	applyPay := &service.ApplyPayOrderInfo{PayOrder: payOrder, OrderList: make([]*service.ApplyOrderInfo, 0)}
	for _, orderInfo := range req.OrderList {
		applyInfo := &service.ApplyOrderInfo{}
		orderDao := conv.ConvertToOrderDao(uid, wxUser.PhoneNumber, orderInfo.ID, req.BuildingID, req.Floor, req.Room)
		if orderDao == nil {
			logger.Warn(orderServerLogTag, "Convert OrderDao Failed|Req:%#v", *req)
			continue
		}
		orderItems := conv.ConvertToOrderDetailDao(orderInfo.OrderItems)

		applyInfo.Order = orderDao
		applyInfo.Items = orderItems
		applyPay.OrderList = append(applyPay.OrderList, applyInfo)
	}

	if len(applyPay.OrderList) == 0 {
		logger.Warn(orderServerLogTag, "ApplyList Length Zero")
		return "", enum.ParamsError, "ID不合法"
	}

	applyPay.PayOrder.MealTime = applyPay.OrderList[0].Order.OrderDate
	prepareID, totalAmount, payAmount, err := os.orderService.ApplyPayOrder(applyPay, dishMap, discountLevel)
	if err != nil {
		logger.Warn(orderServerLogTag, "ApplyPayOrder Failed|Err:%v", err)
		return "", enum.SqlError, err.Error()
	}

	req.TotalAmount = totalAmount
	if payMethod == enum.PayMethodWeChat {
		req.PaymentAmount = payAmount
	}
	return prepareID, enum.Success, ""
}

func (os *OrderServer) RequestApplyCashOrder(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ApplyCashOrderReq)
	uid := req.Uid

	prepareID, code, msg := os.ProcessApplyOrder(uid, req.PayOrderInfo, enum.PayMethodCash)
	if code != enum.Success {
		res.Code = code
		res.Msg = msg
		return
	}

	resData := &dto.ApplyOrderRes{
		PayOrderInfo: req.PayOrderInfo,
		PrepareID:    prepareID,
	}
	res.Data = resData
}

func (os *OrderServer) RequestCancelPayOrder(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.CancelPayOrderReq)
	err := os.orderService.CancelPayOrder(req.OrderID)
	if err != nil {
		logger.Warn(orderServerLogTag, "CancelPayOrder Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
}

func (os *OrderServer) RequestDeliverOrder(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.DeliverOrderReq)
	err := os.orderService.DeliverOrder(req.OrderID)
	if err != nil {
		logger.Warn(orderServerLogTag, "RequestDeliverOrder Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
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
			ID:             payOrder.ID,
			OrderList:      make([]*dto.OrderInfo, 0),
			Floor:          payOrder.Floor,
			Room:           payOrder.Room,
			TotalAmount:    payOrder.TotalAmount,
			PayMethod:      payOrder.PayMethod,
			PaymentAmount:  payOrder.PayAmount,
			DiscountAmount: payOrder.DiscountAmount,
			Status:         payOrder.Status,
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

	orderInfoList := conv.ConvertToOrderInfoList(orderList, detailMap, dishMap)
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

func (os *OrderServer) RequestOrderDishAnalysis(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.OrderDishAnalysisReq)
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	startTime := utils.GetZeroTime(req.OrderDate)
	endTime := utils.GetDayEndTime(req.OrderDate)
	_, detailMap, err := os.orderService.GetAllOrder(req.MealType, startTime, endTime, enum.OrderPaid, req.DishType)
	if err != nil {
		logger.Warn(orderServerLogTag, "GetAllOrder Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	orderNumberMap := make(map[uint32]int32, 0)
	quantityMap := make(map[uint32]int32)
	for _, detailList := range detailMap {
		for _, detailInfo := range detailList {
			if _, ok := quantityMap[detailInfo.DishID]; ok {
				quantityMap[detailInfo.DishID] += detailInfo.Quantity
				orderNumberMap[detailInfo.DishID] += 1
			} else {
				quantityMap[detailInfo.DishID] = detailInfo.Quantity
				orderNumberMap[detailInfo.DishID] = 1
			}
		}

	}
	retData := &dto.OrderDishAnalysisRes{Summary: make([]*dto.OrderDishSummaryInfo, 0)}
	for dishID, quantity := range quantityMap {
		summary := &dto.OrderDishSummaryInfo{
			DishID:      dishID,
			DishName:    dishMap[dishID].DishName,
			DishType:    dishMap[dishID].DishType,
			Quantity:    quantity,
			OrderNumber: orderNumberMap[dishID],
		}
		retData.Summary = append(retData.Summary, summary)
	}
	res.Data = retData
}

func (os *OrderServer) RequestFloorFilter(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.FloorFilterReq)

	startTime := utils.GetZeroTime(req.OrderDate)
	endTime := utils.GetDayEndTime(req.OrderDate)
	floorList, err := os.orderService.GetFloors(req.BuildingID, req.OrderStatus,
		startTime, endTime, req.MealType)
	if err != nil {
		logger.Warn(orderServerLogTag, "GetFloors Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	retData := &dto.FloorFilterRes{Floors: floorList}
	res.Data = retData
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
	if req.StartTime > 0 {
		req.StartTime = utils.GetZeroTime(req.StartTime)
	}
	if req.EndTime > 0 {
		req.EndTime = utils.GetDayEndTime(req.EndTime)
	}
	orderList, totalNumber, detailMap, err := os.orderService.GetOrderList(orderIDList, req.Uid, req.MealType,
		req.BuildingID, req.Floor, req.Room, req.OrderStatus, req.Page, req.PageSize, req.StartTime, req.EndTime)
	if err != nil {
		logger.Warn(orderServerLogTag, "GetOrderList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	orderInfoList := conv.ConvertToOrderInfoList(orderList, detailMap, dishMap)
	resData := &dto.OrderListRes{OrderList: orderInfoList,
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: totalNumber,
		},
	}
	res.Data = resData
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

func (os *OrderServer) RequestModifyCart(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyCartReq)

	uid := req.Uid
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	ids := strings.Split(req.ItemID, "_")
	if len(ids) != 4 {
		res.Code = enum.ParamsError
		res.Msg = "id不合法"
		return
	}

	_, cartDetails, err := os.orderService.ModifyCart(uid, req.ItemID, req.Quantity)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	retData := &dto.ModifyCartRes{
		GoodsMap:   make(map[string]float64),
		TotalCost:  0.0,
		TotalGoods: 0,
	}
	for _, detail := range cartDetails {
		retData.GoodsMap[detail.ItemID] = detail.Quantity
		dishID, _ := conv.ConvertDishID(detail.ItemID)
		dish, ok := dishMap[dishID]
		if ok {
			retData.TotalCost += dish.Price * detail.Quantity
		}
		retData.TotalGoods += detail.Quantity
	}
	res.Data = retData
}
