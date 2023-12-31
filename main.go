package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/canteen_management/config"
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/server"
	"github.com/canteen_management/utils"
	"github.com/gin-gonic/gin"
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
	err := initTokenService()
	if err != nil {
		logger.Warn(serverLogTag, "InitTokenService Failed|Err:%v", err)
		return
	}
	router.Use(CheckSign)
	//router.Use(CheckToken)

	err = HandleAuthApi(router)
	if err != nil {
		logger.Warn(serverLogTag, "HandleAuthApi Failed|Err:%v", err)
		return
	}
	err = HandleUserApi(router)
	if err != nil {
		logger.Warn(serverLogTag, "HandleUserApi Failed|Err:%v", err)
		return
	}
	err = HandlePurchaseApi(router)
	if err != nil {
		logger.Warn(serverLogTag, "HandlePurchaseApi Failed|Err:%v", err)
		return
	}
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
	err = HandleStatisticApi(router)
	if err != nil {
		logger.Warn(serverLogTag, "HandleStatisticApi Failed|Err:%v", err)
		return
	}
	err = HandleOrderApi(router)
	if err != nil {
		logger.Warn(serverLogTag, "HandleOrderApi Failed|Err:%v", err)
		return
	}
	HandleUploadApi(router)

	router.Run(":8081")
}

func HandleAuthApi(router *gin.Engine) error {
	authRouter := router.Group("/api/auth")
	userServer, err := server.NewUserServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewUserServer Failed|Err:%v", err)
		return err
	}
	authRouter.POST("/adminLogin", NewHandler(userServer.RequestAdminLogin,
		func() interface{} { return new(dto.AdminLoginReq) }))
	authRouter.POST("/canteenLogin", NewHandler(userServer.RequestCanteenLogin,
		func() interface{} { return new(dto.CanteenLoginReq) }))
	authRouter.POST("/kitchenLogin", NewHandler(userServer.RequestKitchenLogin,
		func() interface{} { return new(dto.KitchenLoginReq) }))
	return nil
}

func HandleUserApi(router *gin.Engine) error {
	userRouter := router.Group("/api/user")
	userServer, err := server.NewUserServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewUserServer Failed|Err:%v", err)
		return err
	}
	userRouter.Use(CheckToken)
	userRouter.POST("/routerTypeList", NewHandler(userServer.RequestRouterTypeList,
		func() interface{} { return new(dto.RouterTypeListReq) }))
	userRouter.POST("/modifyRouterType", NewHandler(userServer.RequestModifyRouterType,
		func() interface{} { return new(dto.ModifyRouterTypeReq) }))
	userRouter.POST("/routerList", NewHandler(userServer.RequestRouterList,
		func() interface{} { return new(dto.RouterListReq) }))
	userRouter.POST("/modifyRouter", NewHandler(userServer.RequestModifyRouter,
		func() interface{} { return new(dto.ModifyRouterReq) }))
	userRouter.POST("/bindPhoneNumber", NewHandler(userServer.RequestBindPhoneNumber,
		func() interface{} { return new(dto.BindPhoneNumberReq) }))
	userRouter.POST("/canteenUserCenter", NewHandler(userServer.RequestCanteenUserCenter,
		func() interface{} { return new(dto.CanteenUserCenterReq) }))
	userRouter.POST("/kitchenUserCenter", NewHandler(userServer.RequestKitchenUserCenter,
		func() interface{} { return new(dto.KitchenUserCenterReq) }))

	userRouter.POST("/orderUserList", NewHandler(userServer.RequestOrderUserList,
		func() interface{} { return new(dto.OrderUserListReq) }))
	userRouter.POST("/modifyOrderUser", NewHandler(userServer.RequestModifyOrderUser,
		func() interface{} { return new(dto.ModifyOrderUserReq) }))
	userRouter.POST("/bindOrderUser", NewHandler(userServer.RequestBindOrderUser,
		func() interface{} { return new(dto.BindOrderUserReq) }))

	userRouter.POST("/adminUserList", NewHandler(userServer.RequestAdminUserList,
		func() interface{} { return new(dto.AdminUserListReq) }))
	userRouter.POST("/modifyAdminUser", NewHandler(userServer.RequestModifyAdminUser,
		func() interface{} { return new(dto.ModifyAdminUserReq) }))
	userRouter.POST("/bindAdminUser", NewHandler(userServer.RequestBindAdminUser,
		func() interface{} { return new(dto.BindAdminReq) }))
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
	purchaseRouter.Use(CheckToken)
	purchaseRouter.POST("/purchaseList", NewHandler(purchaseServer.RequestPurchaseList,
		func() interface{} { return new(dto.PurchaseListReq) }))
	purchaseRouter.POST("/applyPurchase", NewHandler(purchaseServer.RequestApplyPurchase,
		func() interface{} { return new(dto.ApplyPurchaseReq) }))
	purchaseRouter.POST("/reviewPurchase", NewHandler(purchaseServer.RequestReviewPurchase,
		func() interface{} { return new(dto.ReviewPurchaseReq) }))
	purchaseRouter.POST("/confirmPurchase", NewHandler(purchaseServer.RequestConfirmPurchase,
		func() interface{} { return new(dto.ConfirmPurchaseReq) }))
	purchaseRouter.POST("/receivePurchase", NewHandler(purchaseServer.RequestReceivePurchase,
		func() interface{} { return new(dto.ReceivePurchaseReq) }))

	purchaseRouter.POST("/applyOutbound", NewHandler(purchaseServer.RequestApplyOutbound,
		func() interface{} { return new(dto.ApplyOutboundReq) }))
	purchaseRouter.POST("/reviewOutbound", NewHandler(purchaseServer.RequestReviewOutbound,
		func() interface{} { return new(dto.ReviewOutboundReq) }))
	purchaseRouter.POST("/finishOutbound", NewHandler(purchaseServer.RequestFinishOutbound,
		func() interface{} { return new(dto.FinishOutboundReq) }))
	purchaseRouter.POST("/outboundList", NewHandler(purchaseServer.RequestOutboundOrderList,
		func() interface{} { return new(dto.OutboundListReq) }))

	purchaseRouter.POST("/supplierList", NewHandler(purchaseServer.RequestSupplierList,
		func() interface{} { return new(dto.SupplierListReq) }))
	purchaseRouter.POST("/modifySupplier", NewHandler(purchaseServer.RequestModifySupplier,
		func() interface{} { return new(dto.ModifySupplierReq) }))
	purchaseRouter.POST("/bindSupplier", NewHandler(purchaseServer.RequestBindSupplier,
		func() interface{} { return new(dto.BindSupplierReq) }))
	purchaseRouter.POST("/renewSupplier", NewHandler(purchaseServer.RequestRenewSupplier,
		func() interface{} { return new(dto.RenewSupplierReq) }))
	return nil
}

func HandleStorehouseApi(router *gin.Engine) error {
	storeRouter := router.Group("/api/store")
	storeServer, err := server.NewStorehouseServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewStorehouseServer Failed|Err:%v", err)
		return err
	}
	storeRouter.Use(CheckToken)
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
	storeRouter.POST("/goodsNodeList", NewHandler(storeServer.RequestGoodsNodeList,
		func() interface{} { return new(dto.GoodsNodeListReq) }))
	storeRouter.POST("/modifyGoods", NewHandler(storeServer.RequestModifyGoods,
		func() interface{} { return new(dto.ModifyGoodsInfoReq) }))
	storeRouter.POST("/goodsHistory", NewHandler(storeServer.RequestGoodsHistory,
		func() interface{} { return new(dto.GoodsHistoryReq) }))

	storeRouter.POST("/goodsPriceList", NewHandler(storeServer.RequestGoodsPriceList,
		func() interface{} { return new(dto.GoodsPriceListReq) }))
	storeRouter.POST("/modifyGoodsPrice", NewHandler(storeServer.RequestModifyGoodsPrice,
		func() interface{} { return new(dto.ModifyGoodsPriceReq) }))

	storeRouter.POST("/inventoryList", NewHandler(storeServer.RequestInventoryOrderList,
		func() interface{} { return new(dto.InventoryListReq) }))
	storeRouter.POST("/startInventory", NewHandler(storeServer.RequestStartInventory,
		func() interface{} { return new(dto.StartInventoryReq) }))
	storeRouter.POST("/updateInventory", NewHandler(storeServer.RequestInventory,
		func() interface{} { return new(dto.InventoryReq) }))
	storeRouter.POST("/applyInventory", NewHandler(storeServer.RequestApplyInventory,
		func() interface{} { return new(dto.ApplyInventoryReq) }))
	storeRouter.POST("/confirmInventory", NewHandler(storeServer.RequestConfirmInventory,
		func() interface{} { return new(dto.ConfirmInventoryReq) }))
	storeRouter.POST("/reviewInventory", NewHandler(storeServer.RequestReviewInventory,
		func() interface{} { return new(dto.ReviewInventoryReq) }))
	return nil
}

func HandleMenuApi(router *gin.Engine) error {
	menuRouter := router.Group("/api/menu")
	menuServer, err := server.NewMenuServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewMenuServer Failed|Err:%v", err)
		return err
	}
	menuRouter.Use(CheckToken)
	menuRouter.POST("/dishTypeList", NewHandler(menuServer.RequestDishTypeList,
		func() interface{} { return new(dto.DishTypeListReq) }))
	menuRouter.POST("/modifyDishType", NewHandler(menuServer.RequestModifyDishType,
		func() interface{} { return new(dto.ModifyDishTypeReq) }))

	menuRouter.POST("/dishList", NewHandler(menuServer.RequestDishList,
		func() interface{} { return new(dto.DishListReq) }))
	menuRouter.POST("/modifyDish", NewHandler(menuServer.RequestModifyDish,
		func() interface{} { return new(dto.ModifyDishReq) }))
	menuRouter.POST("/batchModifyDish", NewHandler(menuServer.RequestBatchModifyDish,
		func() interface{} { return new(dto.BatchModifyDishReq) }))

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

func HandleStatisticApi(router *gin.Engine) error {
	statisticRouter := router.Group("/api/statistic")
	statisticServer, err := server.NewStatisticServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewOrderServer Failed|Err:%v", err)
		return err
	}
	statisticRouter.Use(CheckToken)
	statisticRouter.POST("dashboard", NewHandler(statisticServer.RequestDashboard,
		func() interface{} { return new(dto.OrderMenuReq) }))
	return nil
}

func HandleOrderApi(router *gin.Engine) error {
	orderRouter := router.Group("/api/order")
	orderServer, err := server.NewOrderServer(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewOrderServer Failed|Err:%v", err)
		return err
	}
	orderRouter.Use(CheckToken)
	orderRouter.POST("/orderMenuList", NewHandler(orderServer.RequestOrderMenu,
		func() interface{} { return new(dto.OrderMenuReq) }))
	orderRouter.POST("/orderAnalysis", NewHandler(orderServer.RequestOrderDishAnalysis,
		func() interface{} { return new(dto.OrderDishAnalysisReq) }))
	orderRouter.POST("/applyPayOrder", NewHandler(orderServer.RequestApplyOrder,
		func() interface{} { return new(dto.ApplyPayOrderReq) }))
	orderRouter.POST("/applyCashOrder", NewHandler(orderServer.RequestApplyCashOrder,
		func() interface{} { return new(dto.ApplyCashOrderReq) }))
	orderRouter.POST("/applyStaffOrder", NewHandler(orderServer.RequestApplyStaffOrder,
		func() interface{} { return new(dto.ApplyCashOrderReq) }))
	orderRouter.POST("/cancelPayOrder", NewHandler(orderServer.RequestCancelPayOrder,
		func() interface{} { return new(dto.CancelPayOrderReq) }))
	orderRouter.POST("/finishPayOrder", NewHandler(orderServer.RequestFinishPayOrder,
		func() interface{} { return new(dto.FinishPayOrderReq) }))
	orderRouter.POST("/cancelCashOrder", NewHandler(orderServer.RequestCancelCashOrder,
		func() interface{} { return new(dto.CancelPayOrderReq) }))
	orderRouter.POST("/finishCashOrder", NewHandler(orderServer.RequestFinishCashOrder,
		func() interface{} { return new(dto.FinishPayOrderReq) }))
	orderRouter.POST("/payOrderList", NewHandler(orderServer.RequestPayOrderList,
		func() interface{} { return new(dto.PayOrderListReq) }))
	orderRouter.POST("/orderList", NewHandler(orderServer.RequestOrderList,
		func() interface{} { return new(dto.OrderListReq) }))

	orderRouter.POST("/floorFilter", NewHandler(orderServer.RequestFloorFilter,
		func() interface{} { return new(dto.FloorFilterReq) }))
	orderRouter.POST("/deliverOrder", NewHandler(orderServer.RequestDeliverOrder,
		func() interface{} { return new(dto.DeliverOrderReq) }))
	orderRouter.POST("/modifyCart", NewHandler(orderServer.RequestModifyCart,
		func() interface{} { return new(dto.ModifyCartReq) }))

	orderRouter.POST("/orderDiscountList", NewHandler(orderServer.RequestDiscountList,
		func() interface{} { return new(dto.OrderDiscountListReq) }))
	orderRouter.POST("/modifyOrderDiscount", NewHandler(orderServer.RequestModifyDiscount,
		func() interface{} { return new(dto.ModifyOrderDiscountReq) }))
	return nil
}

func HandleUploadApi(router *gin.Engine) {
	uploadRouter := router.Group("/api/")

	uploadServer := server.NewUploadServer()
	uploadRouter.POST("upload", NewHandler(uploadServer.RequestUpload,
		func() interface{} { return new(dto.Request) }))
	uploadRouter.POST("uploadBase64", NewHandler(uploadServer.RequestUploadBase64,
		func() interface{} { return new(dto.UploadBase64Req) }))
}

type RequestDealFunc func(*gin.Context, interface{}, *dto.Response)
type ReqGenerateFunc func() interface{}

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
		logger.Debug(serverLogTag, "RawReq:%#v", req)
		if pageReq, ok := req.(dto.PaginationQ); ok {
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
			logger.Debug(serverLogTag, "Path:%v|Res:%+v|%+v", ctx.FullPath(), res.Code, res.Data)
			ctx.JSON(http.StatusOK, res)
		}()
		dealFunc(ctx, req, res)
		if res.Msg == "" {
			res.Msg = enum.GetMessage(res.Code)
		}
		if res.Code != enum.Success {
			res.Success = false
		}
		if pageRes, ok := res.Data.(dto.PaginationS); ok {
			pageRes.Format()
		}
	}
	return handleFunc
}

var tokenModel *model.TokenModel
var tokenTimeOut = time.Minute * 30
var updateTokenTimeOut = time.Minute * 5

func initTokenService() error {
	sqlCli, err := utils.NewMysqlClient(config.Config.MysqlConfig)
	if err != nil {
		logger.Warn(serverLogTag, "NewMysqlClient Failed|Err:%v", err)
		return err
	}
	tokenModel = model.NewTokenModelWithDB(sqlCli)

	logger.Info(serverLogTag, "tokenTimeOut:%v|updateTokenTimeOut:%v", tokenTimeOut, updateTokenTimeOut)
	return nil
}

func CheckTokenImpl(token string) (*model.TokenDAO, enum.ErrorCode) {
	if token == "" {
		return nil, enum.TokenCheckFailed
	}
	tokenDAO := tokenModel.Get(token)
	if tokenDAO == nil {
		return nil, enum.TokenCheckFailed
	}
	nowTime := time.Now()
	if tokenDAO.UpdateAt.Add(tokenTimeOut).Before(nowTime) {
		return nil, enum.TokenTimeout
	}

	if tokenDAO.UpdateAt.Add(updateTokenTimeOut).Before(nowTime) {
		tokenDAO.UpdateAt = nowTime
		tokenModel.UpdateTokenWithTx(nil, tokenDAO, "updated_at")
	}
	return tokenDAO, enum.Success
}

func CheckToken(c *gin.Context) {
	custom := dto.GetCustomContextInfo(c)
	token, ok := custom.ParamMap[config.TokenKey]
	if !ok {
		logger.Warn(serverLogTag, "Get Token Failed")
		c.AbortWithStatusJSON(http.StatusOK, dto.Response{Code: enum.TokenCheckFailed, Msg: "token not found"})
		return
	}

	tokenDao, code := CheckTokenImpl(token)
	if code != enum.Success {
		logger.Warn(serverLogTag, "Check Token Failed|Token:%v|Code:%v", token, code)
		c.AbortWithStatusJSON(http.StatusOK, dto.Response{Code: enum.TokenCheckFailed, Msg: "token check failed"})
		return
	}

	custom.Token = tokenDao
	c.Set(config.CustomKey, custom)
	c.Next()
}

func CheckSign(c *gin.Context) {
	paramJson := make(map[string]json.RawMessage)
	jsonStr, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Warn(serverLogTag, "Read RequestBody Failed|Err:%v")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(jsonStr, &paramJson)
	if err != nil {
		logger.Warn(serverLogTag, "Parse Request Json Failed|Err:%v")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	paramMap := make(map[string]string)
	if signStr, ok := paramJson["sign"]; ok {
		for k, v := range paramJson {
			paramMap[k] = strings.Trim(string(v), "\"")
		}
		expectSign := utils.GenerateSign(paramMap, config.Config.Secret)
		if expectSign != strings.Trim(string(signStr), "\"") {
			logger.Warn(serverLogTag, "CheckSign Failed|ReqSign:%v|ExpectSign:%v", string(signStr), expectSign)
			//c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CheckSign Failed"))
		} else {
			logger.Debug(serverLogTag, "check sign ok")
		}
	} else {
		logger.Warn(serverLogTag, "no sign")
		//c.AbortWithError(http.StatusBadRequest, fmt.Errorf("no sign"))
		//return
		for k, v := range paramJson {
			paramMap[k] = strings.Trim(string(v), "\"")
		}
	}

	logger.Debug(serverLogTag, "Params:v", paramMap)

	custom := dto.GetCustomContextInfo(c)
	custom.ParamMap = paramMap
	c.Set(config.CustomKey, custom)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonStr))
	c.Next()
}
