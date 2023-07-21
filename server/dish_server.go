package server

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/service"
	"github.com/canteen_management/utils"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	dishServerLogTag = "DishServer"
)

type MenuServer struct {
	dishService *service.DishService
	menuService *service.MenuService
}

func NewMenuServer(dbConf utils.Config) (*MenuServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(dishServerLogTag, "NewMenuServer Failed|Err:%v", err)
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
	return &MenuServer{
		dishService: dishService,
		menuService: menuService,
	}, nil
}

func (ms *MenuServer) RequestDishTypeList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	list, err := ms.dishService.GetDishTypeList()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.DishTypeListRes{
		List: ConvertToDishTypeInfoList(list),
	}
}

func (ms *MenuServer) RequestModifyDishType(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyDishTypeReq)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ms.dishService.AddDishType(ConvertFromDishTypeInfo(&req.TypeInfo))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ms.dishService.ModifyDishType(ConvertFromDishTypeInfo(&req.TypeInfo))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(dishServerLogTag, "RequestModifyDishType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
	return
}

func (ms *MenuServer) RequestDishList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.DishListReq)

	dishList, err := ms.dishService.GetDishList(req.DishType)
	if err != nil {
		res.Code = enum.SqlError
		return
	}
	res.Data = &dto.DishListRes{
		DishList: ConvertToDishInfoList(dishList),
	}
}

func (ms *MenuServer) RequestModifyDish(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyDishReq)

	// todo 校验type
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ms.dishService.AddDish(ConvertFromDishInfo(&req.DishInfo))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ms.dishService.ModifyDish(ConvertFromDishInfo(&req.DishInfo))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(dishServerLogTag, "RequestModifyDish Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
	return
}

func (ms *MenuServer) RequestMenuList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.MenuListReq)

	menuList, err := ms.menuService.GetMenuList(req.MenuType, req.TimeStart, req.TimeEnd)
	if err != nil {
		res.Code = enum.SqlError
		return
	}
	res.Data = &dto.MenuListRes{
		MenuList: ConvertToMenuInfoList(menuList, ms.dishService.GetDishIDMap()),
	}
}

func (ms *MenuServer) RequestModifyMenu(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyMenuReq)
	menuDao, err := ConvertFromMenuInfo(&req.Menu)
	if err != nil {
		res.Code = enum.ParamsError
		return
	}
	switch req.Operate {
	case enum.OperateTypeAdd:
		err = ms.menuService.AddMenu(menuDao)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err = ms.menuService.UpdateMenu(menuDao)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(dishServerLogTag, "RequestModifyMenu Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ms *MenuServer) RequestMenuTypeList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	menuTypeList, err := ms.menuService.GetMenuTypeList()
	if err != nil {
		res.Code = enum.SqlError
		return
	}
	res.Data = &dto.MenuTypeListRes{
		TypeList: ConvertToMenuTypeList(menuTypeList),
	}
}

func (ms *MenuServer) RequestModifyMenuType(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyMenuTypeReq)

	menuTypeDao, err := ConvertFromMenuTypeInfo(&req.MenuType)
	if err != nil {
		res.Code = enum.ParamsError
		return
	}
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ms.menuService.AddMenuType(menuTypeDao)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err = ms.menuService.UpdateMenuType(menuTypeDao)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(dishServerLogTag, "RequestModifyMenuType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ms *MenuServer) RequestGenerateMenu(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GenerateMenuReq)

	typeConf := ms.menuService.GetMenuTypeConfig(req.MenuType)
	if typeConf == nil {
		logger.Warn(dishServerLogTag, "MenuType Not Found|MenuType:%v", req.MenuType)
		res.Code = enum.ParamsError
		return
	}
	confMap, err := ConvertFromMenuTypeConfig(typeConf.MenuConfig)
	if err != nil {
		logger.Warn(dishServerLogTag, "MenuTypeConfig Convert Failed|MenuType:%v|Err:%v", req.MenuType, err)
		res.Code = enum.ParamsError
		return
	}

	startTime := time.Unix(req.TimeStart, 0)
	start := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0,
		startTime.Location()).Unix()
	menuList := make([]*dto.MenuInfo, 0)
	for ; start <= req.TimeEnd; start += 3600 * 24 {
		mealList := make([]*dto.MealInfo, 0)
		for _, conf := range confMap {
			totalDishList := make([]*model.Dish, 0)
			for dishType, dishNum := range conf.DishNumberMap {
				dishList := ms.dishService.RandDishByType(dishType, int(dishNum))
				totalDishList = append(totalDishList, dishList...)
			}
			dishList := ConvertToDishInfoList(totalDishList)
			mealInfo := &dto.MealInfo{
				MealName: conf.MealName,
				MealType: conf.MealType,
				DishList: dishList,
			}
			mealList = append(mealList, mealInfo)
		}
		menuList = append(menuList, &dto.MenuInfo{MenuType: req.MenuType, MenuDate: start, MealList: mealList})
	}

	res.Data = &dto.GenerateMenuRes{MenuList: menuList}
}
