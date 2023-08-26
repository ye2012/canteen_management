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
				retDish := &dto.OrderNode{ID: fmt.Sprintf("%v", index), Name: dish.DishName, OrderNodeMap: make(map[string]interface{})}
				retDish.OrderNodeMap[KeyPrice] = dish.Price
				retListByType.Children = append(retListByType.Children, retDish)
			}
			retMeal.Children = append(retMeal.Children, retListByType)
		}

		retData = append(retData, retMeal)
	}
	return retData
}

func (os *OrderServer) RequestApplyOrder(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ApplyOrderReq)
	orderDao := ConvertToOrderDao(0, req.ID, req.Address, req.PickUpMethod)
	if orderDao == nil {
		logger.Warn(orderServerLogTag, "Convert OrderDao Failed|Req:%#v", *req)
		res.Code = enum.ParamsError
		return
	}
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	orderItems := ConvertToOrderDetailDao(req.OrderItems)

	prepareID, err := os.orderService.ApplyOrder(orderDao, orderItems, dishMap, 1)
	if err != nil {
		logger.Warn(orderServerLogTag, "ApplyOrder Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	req.TotalAmount = orderDao.TotalAmount
	req.PaymentAmount = orderDao.PayAmount

	resData := &dto.ApplyOrderRes{
		Order:     req,
		PrepareID: prepareID,
	}
	res.Data = resData
}

func ConvertToOrderDao(uid uint32, ID, addr string, pickUpMethod uint8) *model.OrderDao {
	ids := strings.Split(ID, "_")
	if len(ids) != 2 {
		logger.Warn(orderServerLogTag, "ID illegal|ID:%v", ID)
		return nil
	}
	mealTime, _ := strconv.ParseInt(ids[0], 10, 32)
	mealType, _ := strconv.ParseInt(ids[1], 10, 32)

	return &model.OrderDao{
		OrderDate:    time.Unix(mealTime, 0),
		MealType:     uint8(mealType),
		Uid:          uid,
		PickUpMethod: pickUpMethod,
		Address:      addr,
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

func (os *OrderServer) RequestOrderList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.OrderListReq)
	dishMap, err := os.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(orderServerLogTag, "GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}
	orderList, detailMap, err := os.orderService.GetOrderList(req.OrderID, req.Uid, req.Page, req.PageSize, req.OrderStatus)
	if err != nil {
		logger.Warn(orderServerLogTag, "GetOrderList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	orderInfoList := ConvertToOrderInfoList(orderList, detailMap, dishMap)
	res.Data = orderInfoList
}

func ConvertToOrderInfoList(orderList []*model.OrderDao, detailMap map[uint32][]*model.OrderDetail,
	dishMap map[uint32]*model.Dish) []*dto.OrderInfo {
	retList := make([]*dto.OrderInfo, 0)
	for _, order := range orderList {
		retInfo := &dto.OrderInfo{
			ID:            fmt.Sprintf("%v_%v", order.OrderDate.Unix(), order.MealType),
			UnionID:       fmt.Sprintf("%v", order.UnionID),
			OrderID:       fmt.Sprintf("%v", order.ID),
			OrderNo:       "",
			Address:       order.Address,
			PickUpMethod:  order.PickUpMethod,
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

func (os *OrderServer) RequestOrderUserList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.OrderUserListReq)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize > 1000 {
		req.PageSize = 100
	}
	userList, userNumber, err := os.orderService.GetOrderUserList(req.PhoneNumber, req.DiscountLevel, req.Page, req.PageSize)
	if err != nil {
		logger.Warn(orderServerLogTag, "GetOrderUserList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	userInfoList := make([]*dto.OrderUserInfo, 0, len(userList))
	for _, user := range userList {
		userInfo := &dto.OrderUserInfo{
			ID:            user.ID,
			PhoneNumber:   user.PhoneNumber,
			DiscountLevel: user.DiscountLevel,
		}
		userInfoList = append(userInfoList, userInfo)
	}

	extraPage := uint32(1)
	if userNumber%req.PageSize == 0 {
		extraPage = 0
	}
	retData := &dto.OrderUserListRes{
		UserList:  userInfoList,
		TotalPage: userNumber/req.PageSize + extraPage,
		PageSize:  req.PageSize,
		Page:      req.Page,
	}
	res.Data = retData
}

func (os *OrderServer) RequestModifyOrderUser(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyOrderUserReq)
	if len(req.UserList) == 0 {
		res.Code = enum.ParamsError
		return
	}

	userList := make([]*model.OrderUser, 0, len(req.UserList))
	for _, userInfo := range req.UserList {
		user := &model.OrderUser{
			ID:            userInfo.ID,
			PhoneNumber:   userInfo.PhoneNumber,
			DiscountLevel: userInfo.DiscountLevel,
		}
		userList = append(userList, user)
	}
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := os.orderService.AddOrderUser(userList)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := os.orderService.UpdateOrderUser(userList[0])
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(orderServerLogTag, "RequestModifyOrderUser Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}
