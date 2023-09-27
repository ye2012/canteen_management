package main

import (
	"flag"
	"net/http"
	"reflect"

	"github.com/canteen_management/config"
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/server"
	"github.com/gin-gonic/gin"
	"gopkg.in/gotsunami/coquelicot.v1"
)

const (
	serverLogTag = "Server"
)

var (
	configPath = flag.String("config", "./config.json", "config file abs path")
)

func main() {
	flag.Parsed()
	config.LoadConfig(*configPath)
	StartServer()
}

func StartServer() {
	router := gin.Default()

	err := HandleUserApi(router)
	if err != nil {
		logger.Warn(serverLogTag, "HandleUserApi Failed|Err:%v", err)
		return
	}
	HandlePurchaseApi(router)
	err = HandleStorehouseApi(router)
	if err != nil {
		logger.Warn(serverLogTag, "HandleStorehouseApi Failed|Err:%v", err)
		return
	}
	err = HandleMenuApi(router)
	if err != nil {
		logger.Warn(serverLogTag, "HandleMenuApi Failed|Err:%v", err)
		return
	}
	HandleStatisticApi(router)
	HandleOrderApi(router)
	HandleUploadApi(router)

	router.Run(":8081")
}

func HandleUserApi(router *gin.Engine) error {
	userRouter := router.Group("/api/user")
	userServer, err := server.NewUserServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewUserServer Failed|Err:%v", err)
		return err
	}
	userRouter.POST("/canteenLogin", NewHandler(userServer.RequestCanteenLogin,
		func() interface{} { return new(dto.CanteenLoginReq) }))
	userRouter.POST("/bindPhoneNumber", NewHandler(userServer.RequestBindPhoneNumber,
		func() interface{} { return new(dto.BindPhoneNumberReq) }))

	userRouter.POST("/orderUserList", NewHandler(userServer.RequestOrderUserList,
		func() interface{} { return new(dto.OrderUserListReq) }))
	userRouter.POST("/modifyOrderUser", NewHandler(userServer.RequestModifyOrderUser,
		func() interface{} { return new(dto.ModifyOrderUserReq) }))
	//userRouter.POST("/modifyUser", NewHandler(storeServer.,
	//	&dto.{}))
	return nil
}

func HandlePurchaseApi(router *gin.Engine) error {
	purchaseRouter := router.Group("/api/purchase")
	purchaseServer, err := server.NewPurchaseServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewPurchaseServer Failed|Err:%v", err)
		return err
	}

	//purchaseRouter.POST("/purchaseList", NewHandler(storeServer.,
	//	&dto.{}))
	//
	//purchaseRouter.POST("/applyPurchase", NewHandler(storeServer.,
	//	&dto.{}))
	//purchaseRouter.POST("/reviewPurchase", NewHandler(storeServer.,
	//	&dto.{}))
	//purchaseRouter.POST("/acceptPurchase", NewHandler(storeServer.,
	//	&dto.{}))
	//purchaseRouter.POST("/finishPurchase", NewHandler(storeServer.,
	//	&dto.{}))
	//
	//
	purchaseRouter.POST("/supplierList", NewHandler(purchaseServer.RequestSupplierList,
		func() interface{} { return new(dto.SupplierListReq) }))
	purchaseRouter.POST("/modifySupplier", NewHandler(purchaseServer.RequestModifySupplier,
		func() interface{} { return new(dto.ModifySupplierReq) }))
	purchaseRouter.POST("/bindSupplier", NewHandler(purchaseServer.RequestBindSupplier,
		func() interface{} { return new(dto.BindSupplierReq) }))
	return nil
}

func HandleStorehouseApi(router *gin.Engine) error {
	storeRouter := router.Group("/api/store")
	storeServer, err := server.NewStorehouseServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewStorehouseServer Failed|Err:%v", err)
		return err
	}

	storeRouter.POST("/storeTypeList", NewHandler(storeServer.RequestStoreTypeList,
		func() interface{} { return new(dto.StoreTypeListReq) }))
	storeRouter.POST("/modifyStoreType", NewHandler(storeServer.RequestModifyStoreType,
		func() interface{} { return new(dto.ModifyStoreTypeReq) }))

	storeRouter.POST("/goodsTypeList", NewHandler(storeServer.RequestGoodsTypeList,
		func() interface{} { return new(dto.GoodsTypeListReq) }))
	storeRouter.POST("/modifyGoodsType", NewHandler(storeServer.RequestModifyGoodsType,
		func() interface{} { return new(dto.ModifyGoodsTypeReq) }))

	storeRouter.POST("/goodsList", NewHandler(storeServer.RequestGoodsList,
		func() interface{} { return new(dto.GoodsListReq) }))
	storeRouter.POST("/modifyGoods", NewHandler(storeServer.RequestModifyGoods,
		func() interface{} { return new(dto.ModifyGoodsInfoReq) }))

	storeRouter.POST("/goodsPriceList", NewHandler(storeServer.RequestGoodsPriceList,
		func() interface{} { return new(dto.GoodsPriceListReq) }))
	storeRouter.POST("/modifyGoodsPrice", NewHandler(storeServer.RequestModifyGoodsPrice,
		func() interface{} { return new(dto.ModifyGoodsPriceReq) }))

	//storeRouter.POST("/applyConsumeGoods", NewHandler(storeServer.,
	//	&dto.{}))
	//storeRouter.POST("/confirmConsumeGoods", NewHandler(storeServer.,
	//	&dto.{}))
	return nil
}

func HandleMenuApi(router *gin.Engine) error {
	menuRouter := router.Group("/api/menu")
	menuServer, err := server.NewMenuServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewMenuServer Failed|Err:%v", err)
		return err
	}
	menuRouter.POST("/dishTypeList", NewHandler(menuServer.RequestDishTypeList,
		func() interface{} { return new(dto.DishTypeListReq) }))
	menuRouter.POST("/modifyDishType", NewHandler(menuServer.RequestModifyDishType,
		func() interface{} { return new(dto.ModifyDishTypeReq) }))

	menuRouter.POST("/dishList", NewHandler(menuServer.RequestDishList,
		func() interface{} { return new(dto.DishListReq) }))
	menuRouter.POST("/modifyDish", NewHandler(menuServer.RequestModifyDish,
		func() interface{} { return new(dto.ModifyDishReq) }))

	menuRouter.POST("/weekMenuList", NewHandler(menuServer.RequestWeekMenuList,
		func() interface{} { return new(dto.WeekMenuListReq) }))
	menuRouter.POST("/weekMenuDetail", NewHandler(menuServer.RequestWeekMenuDetail,
		func() interface{} { return new(dto.WeekMenuDetailReq) }))
	menuRouter.POST("/weekMenuListHead", NewHandler(menuServer.RequestWeekMenuListHead,
		func() interface{} { return new(dto.WeekMenuListReq) }))
	menuRouter.POST("/weekMenuListData", NewHandler(menuServer.RequestWeekMenuListData,
		func() interface{} { return new(dto.WeekMenuListReq) }))
	menuRouter.POST("/weekMenuDetailTable", NewHandler(menuServer.RequestWeekMenuDetailTable,
		func() interface{} { return new(dto.WeekMenuDetailReq) }))
	menuRouter.POST("/modifyWeekMenu", NewHandler(menuServer.RequestModifyWeekMenu,
		func() interface{} { return new(dto.ModifyWeekMenuDetailReq) }))

	menuRouter.POST("/generateStaffMenu", NewHandler(menuServer.RequestGenerateStaffMenu,
		func() interface{} { return new(dto.GenerateStaffMenuReq) }))
	menuRouter.POST("/generateWeekMenu", NewHandler(menuServer.RequestGenerateWeekMenu,
		func() interface{} { return new(dto.GenerateWeekMenuReq) }))

	menuRouter.POST("/staffMenuListHead", NewHandler(menuServer.RequestStaffMenuListHead,
		func() interface{} { return new(dto.StaffMenuListHeadReq) }))
	menuRouter.POST("/staffMenuListData", NewHandler(menuServer.RequestStaffMenuListData,
		func() interface{} { return new(dto.StaffMenuListDataReq) }))
	menuRouter.POST("/staffMenuHead", NewHandler(menuServer.RequestStaffMenuDetailHead,
		func() interface{} { return new(dto.StaffMenuDetailHeadReq) }))
	menuRouter.POST("/staffMenuData", NewHandler(menuServer.RequestStaffMenuDetailData,
		func() interface{} { return new(dto.StaffMenuDetailDataReq) }))
	menuRouter.POST("/modifyStaffMenu", NewHandler(menuServer.RequestModifyStaffMenuDetail,
		func() interface{} { return new(dto.ModifyStaffMenuDetailReq) }))

	menuRouter.POST("/menuTypeListHead", NewHandler(menuServer.RequestMenuTypeListHead,
		func() interface{} { return new(dto.MenuTypeListHeadReq) }))
	menuRouter.POST("/menuTypeListData", NewHandler(menuServer.RequestMenuTypeListData,
		func() interface{} { return new(dto.MenuTypeListDataReq) }))
	menuRouter.POST("/menuTypeDetailHead", NewHandler(menuServer.RequestMenuTypeDetailHead,
		func() interface{} { return new(dto.MenuTypeDetailHeadReq) }))
	menuRouter.POST("/menuTypeDetailData", NewHandler(menuServer.RequestMenuTypeDetailData,
		func() interface{} { return new(dto.MenuTypeDetailDataReq) }))
	menuRouter.POST("/modifyMenuType", NewHandler(menuServer.RequestModifyMenuType,
		func() interface{} { return new(dto.ModifyMenuTypeReq) }))
	return nil
}

func HandleStatisticApi(router *gin.Engine) {
	//statisticRouter := router.Group("/statistic")
	//statisticServer := server.NewStatisticServer()
	//statisticRouter.POST("enterDinerNumber", NewHandler(statisticServer))
}

func HandleOrderApi(router *gin.Engine) error {
	orderRouter := router.Group("/api/order")
	orderServer, err := server.NewOrderServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewOrderServer Failed|Err:%v", err)
		return err
	}
	orderRouter.POST("/orderMenuList", NewHandler(orderServer.RequestOrderMenu,
		func() interface{} { return new(dto.OrderMenuReq) }))
	orderRouter.POST("/applyPayOrder", NewHandler(orderServer.RequestApplyOrder,
		func() interface{} { return new(dto.ApplyPayOrderReq) }))
	orderRouter.POST("/cancelPayOrder", NewHandler(orderServer.RequestCancelOrder,
		func() interface{} { return new(dto.CancelPayOrderReq) }))
	orderRouter.POST("/payOrderList", NewHandler(orderServer.RequestPayOrderList,
		func() interface{} { return new(dto.PayOrderListReq) }))
	orderRouter.POST("/orderList", NewHandler(orderServer.RequestOrderList,
		func() interface{} { return new(dto.OrderListReq) }))

	orderRouter.POST("/orderDiscountList", NewHandler(orderServer.RequestDiscountList,
		func() interface{} { return new(dto.OrderDiscountListReq) }))
	orderRouter.POST("/modifyOrderDiscount", NewHandler(orderServer.RequestModifyDiscount,
		func() interface{} { return new(dto.ModifyOrderDiscountReq) }))
	return nil
}

func HandleUploadApi(router *gin.Engine) {
	orderRouter := router.Group("/api/")
	store := coquelicot.NewStorage("/home/work/private/canteen/files/")
	orderRouter.POST("/upload", NewUploadHandler(store.UploadHandler))
}

type RequestDealFunc func(*gin.Context, interface{}, *dto.Response)
type ReqGenerateFunc func() interface{}

func NewUploadHandler(uploadHandler func(w http.ResponseWriter, r *http.Request)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uploadHandler(ctx.Writer, ctx.Request)
	}
}

func NewHandler(dealFunc RequestDealFunc, reqGen ReqGenerateFunc) gin.HandlerFunc {
	handleFunc := func(ctx *gin.Context) {
		req := reqGen()
		logger.Debug(serverLogTag, "%v", reflect.TypeOf(req))
		res := dto.GetInitResponse()
		err := ctx.ShouldBind(req)
		if err != nil {
			logger.Warn(serverLogTag, "parse req failed|Err:%v", err)
			res.Code = enum.ParseRequestFailed
			ctx.JSON(http.StatusBadRequest, res)
			return
		}
		logger.Debug(serverLogTag, "Req:%#v", req)
		if pageReq, ok := req.(*dto.PaginationReq); ok {
			pageReq.FixPagination()
		}
		if checker, ok := req.(dto.RequestChecker); ok {
			err = checker.CheckParams()
			if err != nil {
				res.Code = enum.ParamsError
				res.Msg = err.Error()
				ctx.JSON(http.StatusOK, res)
				return
			}
		}
		defer func() {
			logger.Debug(serverLogTag, "Path:%v|Res:%+v", ctx.FullPath(), res)
			ctx.JSON(http.StatusOK, res)
		}()
		dealFunc(ctx, req, res)
		if res.Msg == "" {
			res.Msg = enum.GetMessage(res.Code)
		}
		if res.Code != enum.Success {
			res.Success = false
		}
		if pageRes, ok := res.Data.(*dto.PaginationRes); ok {
			pageRes.Format()
		}
	}
	return handleFunc
}

//func CheckTokenImpl(token string) (*model.TokenDAO, enum.ErrorCode) {
//	if token == "" {
//		return nil, enum.TokenCheckFailed
//	}
//	tokenDAO := tokenModel.Get(token)
//	if tokenDAO == nil {
//		return nil, enum.TokenCheckFailed
//	}
//	nowTime := time.Now()
//	if tokenDAO.MTime.Add(tokenTimeOut).Before(nowTime) {
//		return nil, enum.TokenTimeout
//	}
//
//	if tokenDAO.MTime.Add(updateTokenTimeOut).Before(nowTime) {
//		tokenDAO.MTime = nowTime
//		tokenModel.UpdateModifyTime(tokenDAO)
//	}
//	return tokenDAO, enum.Success
//}

//func CheckToken(c *gin.Context) {
//	custom := dto.GetCustomContextInfo(c)
//	token, ok := custom.ParamMap[config.TokenKey]
//	if !ok {
//		logger.Warn(serverLogTag, "Get Token Failed")
//		c.AbortWithStatusJSON(http.StatusOK, dto.Response{Code: enum.ParamsError, Msg: "token not found"})
//		return
//	}
//
//	tokenDao, code := CheckTokenImpl(token)
//	if code != enum.Success {
//		logger.Warn(serverLogTag, "Check Token Failed|Token:%v|Code:%v", token, code)
//		c.AbortWithStatusJSON(http.StatusOK, dto.Response{Code: enum.ParamsError, Msg: "token check failed"})
//		return
//	}
//
//	custom.Token = tokenDao
//	c.Set(config.CustomKey, custom)
//	c.Next()
//}
