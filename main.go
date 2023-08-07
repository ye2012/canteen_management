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

	HandleUserApi(router)
	HandlePurchaseApi(router)
	err := HandleStorehouseApi(router)
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

	router.Run(":8081")
}

func HandleUserApi(router *gin.Engine) {
	//userRouter := router.Group("/user")
	//
	//userRouter.POST("/login", NewHandler(storeServer.,
	//	&dto.{}))
	//userRouter.POST("/userList", NewHandler(storeServer.,
	//	&dto.{}))
	//userRouter.POST("/modifyUser", NewHandler(storeServer.,
	//	&dto.{}))
}

func HandlePurchaseApi(router *gin.Engine) {
	//purchaseRouter := router.Group("/purchase")
	//
	//purchaseRouter.POST("/purchaseList", NewHandler(storeServer.,
	//	&dto.{}))
	//
	//purchaseRouter.POST("/applyPurchase", NewHandler(storeServer.,
	//	&dto.{}))
	//purchaseRouter.POST("/reviewPurchase", NewHandler(storeServer.,
	//	&dto.{}))
	//purchaseRouter.POST("/confirmPurchase", NewHandler(storeServer.,
	//	&dto.{}))
	//purchaseRouter.POST("/finishPurchase", NewHandler(storeServer.,
	//	&dto.{}))
	//
	//
	//purchaseRouter.POST("/supplierList", NewHandler(storeServer.,
	//	&dto.{}))
	//purchaseRouter.POST("/modifySupplier", NewHandler(storeServer.,
	//	&dto.{}))
	//purchaseRouter.POST("/updateDiscount", NewHandler(storeServer.,
	//	&dto.{}))
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
	menuRouter.POST("/modifyWeekMenu", NewHandler(menuServer.RequestModifyWeekMenu,
		func() interface{} { return new(dto.ModifyWeekMenuReq) }))

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
	menuRouter.POST("/modifyStaffMenu", NewHandler(menuServer.RequestModifyMenuDetail,
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

type RequestDealFunc func(*gin.Context, interface{}, *dto.Response)
type ReqGenerateFunc func() interface{}

func NewHandler(dealFunc RequestDealFunc, reqGen ReqGenerateFunc) gin.HandlerFunc {
	handleFunc := func(ctx *gin.Context) {
		req := reqGen()
		logger.Debug(serverLogTag, "%v", reflect.TypeOf(req))
		res := dto.GetInitResponse()
		defer func() {
			logger.Debug(serverLogTag, "Path:%v|Res:%+v", ctx.FullPath(), res)
			ctx.JSON(http.StatusOK, res)
		}()
		err := ctx.ShouldBind(req)
		if err != nil {
			logger.Warn(serverLogTag, "parse req failed|Err:%v", err)
			res.Code = enum.ParseRequestFailed
			return
		}
		logger.Debug(serverLogTag, "Req:%#v", req)
		dealFunc(ctx, req, res)
		if res.Code != 0 {
			res.Success = false
		}
	}
	return handleFunc
}
