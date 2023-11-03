package server

import (
	"github.com/canteen_management/conv"
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
	menuServerLogTag = "MenuServer"
)

type MenuServer struct {
	dishService *service.DishService
	menuService *service.MenuService
}

func NewMenuServer(dbConf utils.Config) (*MenuServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(menuServerLogTag, "NewMenuServer Failed|Err:%v", err)
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
	req := rawReq.(*dto.DishTypeListReq)

	list, count, err := ms.dishService.GetDishTypeList(req.MasterTypeID, req.IncludeMaserType, req.Page, req.PageSize)
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.DishTypeListRes{
		List: conv.ConvertToDishTypeInfoList(list),
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: count,
		},
	}
}

func (ms *MenuServer) RequestModifyDishType(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyDishTypeReq)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ms.dishService.AddDishType(conv.ConvertFromDishTypeInfo(&req.TypeInfo))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ms.dishService.ModifyDishType(conv.ConvertFromDishTypeInfo(&req.TypeInfo))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeDel:
		err := ms.dishService.DeleteDishType(req.TypeInfo.DishTypeID)
		if err != nil {
			res.Code = enum.SqlError
			res.Msg = err.Error()
			return
		}
	default:
		logger.Warn(menuServerLogTag, "RequestModifyDishType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
	return
}

func (ms *MenuServer) RequestDishList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.DishListReq)
	typeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		res.Code = enum.SystemError
		return
	}

	dishList, dishCount, err := ms.dishService.GetDishList(req.DishType, req.Page, req.PageSize)
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.DishListRes{
		DishList: conv.ConvertToDishInfoList(dishList, typeMap),
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: dishCount,
		},
	}
}

func (ms *MenuServer) RequestModifyDish(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyDishReq)

	logger.Info(menuServerLogTag, "RequestModifyDish Req:%#v", req)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ms.dishService.AddDish(conv.ConvertFromDishInfo(&req.DishInfo))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ms.dishService.ModifyDish(conv.ConvertFromDishInfo(&req.DishInfo))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(menuServerLogTag, "RequestModifyDish Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
	return
}

func (ms *MenuServer) RequestWeekMenuList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.WeekMenuListReq)

	menuList, err := ms.menuService.GetWeekMenuList(req.MenuType, req.TimeStart, req.TimeEnd, req.Page, req.PageSize)
	if err != nil {
		res.Code = enum.SqlError
		return
	}
	dishIDMap, err := ms.dishService.GetDishIDMap()
	if err != nil {
		res.Code = enum.SqlError
		return
	}
	res.Data = &dto.WeekMenuListRes{
		MenuList: conv.ConvertToWeekMenuList(menuList, dishIDMap),
	}
}

func (ms *MenuServer) RequestWeekMenuListHead(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	head := conv.GenerateWeekMenuListTableHead()
	res.Data = head
}

func (ms *MenuServer) RequestWeekMenuListData(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.WeekMenuListHeadReq)
	dishMap, err := ms.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestWeekMenuListData GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	menuList, err := ms.menuService.GetWeekMenuList(req.MenuType, req.TimeStart, req.TimeStart, req.Page, req.PageSize)
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestWeekMenuListData GetMenuList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	res.Data = conv.GenerateWeekMenuListTableData(menuList, dishMap)
}

func (ms *MenuServer) RequestWeekMenuDetailTable(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.WeekMenuDetailTableReq)
	dishMap, err := ms.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestWeekMenuDetailTable GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	dishTypeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestWeekMenuDetailTable GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	menu, err := ms.menuService.GetWeekMenu(req.WeekMenuID)
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestStaffMenuListData GetMenuList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	retData := &dto.WeekMenuDetailTableRes{}
	retData.Head, retData.Data = conv.GenerateWeekMenuDetailTable(menu, dishMap, dishTypeMap)
	res.Data = retData
}

func (ms *MenuServer) RequestWeekMenuDetail(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.WeekMenuDetailReq)
	dishIDMap, err := ms.dishService.GetDishIDMap()
	if err != nil {
		res.Code = enum.SqlError
		return
	}
	dishTypeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	weekMenu, err := ms.menuService.GetWeekMenu(req.WeekMenuID)
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	resData, err := conv.ConvertToWeekMenuDetail(weekMenu, dishIDMap, dishTypeMap)
	if err != nil {
		res.Code = enum.SystemError
		return
	}
	res.Data = resData
}

func (ms *MenuServer) RequestModifyWeekMenu(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyWeekMenuDetailReq)
	weekMenuDao := conv.ParseWeekMenuDetail(req.WeekMenuRows, req.WeekMenuID, req.MenuTypeID, req.MenuDate)
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ms.menuService.AddWeekMenu(weekMenuDao)
		if err != nil {
			res.Code = enum.SqlError
			res.Msg = err.Error()
			return
		}
	case enum.OperateTypeModify:
		err := ms.menuService.UpdateWeekMenu(weekMenuDao)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(menuServerLogTag, "RequestModifyWeekMenu Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ms *MenuServer) RequestModifyMenuType(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyMenuTypeReq)

	if len(req.MenuTypeRows) != 1 {
		res.Code = enum.ParamsError
		return
	}
	menuTypeDao := conv.ParseMenuTypeDetailData(req.MenuTypeRows[0], req.MenuTypeID, req.MenuTypeName)
	logger.Info(menuServerLogTag, "ModifyMenuType:%#v", menuTypeDao)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ms.menuService.AddMenuType(menuTypeDao)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ms.menuService.UpdateMenuType(menuTypeDao)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(menuServerLogTag, "RequestModifyMenuType Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ms *MenuServer) RequestGenerateStaffMenu(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GenerateStaffMenuReq)
	dishMap, err := ms.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestStaffMenuDetailData GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	dishTypeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestGenerateMenu GetDishTypeMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	typeConf, err := ms.menuService.GetMenuType(req.MenuType)
	if err != nil {
		logger.Warn(menuServerLogTag, "GetMenuType Failed|MenuType:%v|Err:%v", req.MenuType, err)
		res.Code = enum.ParamsError
		return
	}
	confMap := typeConf.ToMenuConfig()
	if confMap == nil {
		logger.Warn(menuServerLogTag, "Get MenuTypeConfig Failed|MenuType:%v|Err:%v", req.MenuType, err)
		res.Code = enum.ParamsError
		return
	}

	menu := &model.Menu{MenuDate: time.Now(), MenuTypeID: req.MenuType}
	menuDishMap := make(map[uint8][]uint32)
	for mealType, numberConf := range confMap {
		totalDishList := make([]uint32, 0)
		for dishType, dishNum := range numberConf {
			dishList := ms.dishService.RandDishIDByType(dishType, int(dishNum))
			totalDishList = append(totalDishList, dishList...)
		}
		menuDishMap[mealType] = totalDishList
	}
	menu.FromMenuConfig(menuDishMap)

	resData := &dto.GenerateStaffMenuRes{
		Head: conv.GenerateStaffDetailTableHead(),
		Data: conv.GenerateStaffDetailTableData(menu, dishMap, dishTypeMap),
	}
	res.Data = resData
}

func (ms *MenuServer) RequestGenerateWeekMenu(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.GenerateWeekMenuReq)
	dishMap, err := ms.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestGenerateWeekMenu GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	dishTypeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "GenerateWeekMenu GetDishTypeMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	typeConf, err := ms.menuService.GetMenuType(req.MenuType)
	if err != nil {
		logger.Warn(menuServerLogTag, "GenerateWeekMenu GetMenuType Failed|MenuType:%v|Err:%v", req.MenuType, err)
		res.Code = enum.ParamsError
		return
	}
	confMap := typeConf.ToMenuConfig()
	if confMap == nil {
		logger.Warn(menuServerLogTag, "Get MenuTypeConfig Failed|MenuType:%v|Err:%v", req.MenuType, err)
		res.Code = enum.ParamsError
		return
	}

	startTime := time.Unix(req.TimeStart, 0)
	start := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0,
		startTime.Location()).Unix()
	end := start + 3600*24*7
	menuList := make([]map[uint8][]uint32, 0, 7)
	for ; start < end; start += 3600 * 24 {
		mealMap := make(map[uint8][]uint32)
		for mealType, numberConf := range confMap {
			totalDishList := make([]*model.Dish, 0)
			for dishType, dishNum := range numberConf {
				dishList := ms.dishService.RandDishByType(dishType, int(dishNum))
				totalDishList = append(totalDishList, dishList...)
			}
			dishIDList := make([]uint32, 0, len(totalDishList))
			for _, dish := range totalDishList {
				dishIDList = append(dishIDList, dish.ID)
			}
			mealMap[mealType] = dishIDList
		}
		menuList = append(menuList, mealMap)
	}

	menu := &model.WeekMenu{MenuStartDate: startTime}
	menu.FromWeekMenuConfig(menuList)

	retData := &dto.GenerateWeekMenuRes{}
	retData.Head, retData.Data = conv.GenerateWeekMenuDetailTable(menu, dishMap, dishTypeMap)
	res.Data = retData
}

func (ms *MenuServer) RequestStaffMenuListHead(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	head := conv.GenerateStaffMenuListTableHead()
	res.Data = head
}

func (ms *MenuServer) RequestStaffMenuListData(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.StaffMenuListDataReq)
	dishMap, err := ms.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestStaffMenuDetailData GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	menuList, err := ms.menuService.GetMenuList(2, req.TimeStart, req.TimeStart, req.Page, req.PageSize)
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestStaffMenuListData GetMenuList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	data := conv.GenerateStaffMenuListTableData(menuList, dishMap)
	res.Data = data
}

func (ms *MenuServer) RequestStaffMenuDetailHead(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	head := conv.GenerateStaffDetailTableHead()
	res.Data = head
}

func (ms *MenuServer) RequestStaffMenuDetailData(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.StaffMenuDetailDataReq)
	dishMap, err := ms.dishService.GetDishIDMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestStaffMenuDetailData GetDishIDMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	dishTypeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestStaffMenuDetailData GetDishTypeMap Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	menu, err := ms.menuService.GetMenu(req.StaffMenuID)
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestStaffMenuDetailData GetMenu Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	data := conv.GenerateStaffDetailTableData(menu, dishMap, dishTypeMap)
	res.Data = data
}

func (ms *MenuServer) RequestModifyStaffMenuDetail(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyStaffMenuDetailReq)
	staffMenuDao := conv.ParseStaffMenuDetailData(req.StaffMenuRows, req.StaffMenuID, req.MenuTypeID, req.MenuDate)
	logger.Info(menuServerLogTag, "ModifyMenuDetail:%#v", staffMenuDao)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ms.menuService.AddMenu(staffMenuDao)
		if err != nil {
			res.Code = enum.SqlError
			res.Msg = err.Error()
			return
		}
	case enum.OperateTypeModify:
		err := ms.menuService.UpdateMenu(staffMenuDao)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(menuServerLogTag, "RequestModifyMenuDetail Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ms *MenuServer) RequestMenuTypeListHead(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	head := conv.GenerateMenuTypeListTableHead()
	res.Data = head
}

func (ms *MenuServer) RequestMenuTypeListData(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	dishTypeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	menuTypeList, err := ms.menuService.GetMenuTypeList()
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestMenuTypeListData GetMenuTypeList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}
	data := conv.GenerateMenuTypeListTableData(menuTypeList, dishTypeMap)
	res.Data = data
}

func (ms *MenuServer) RequestMenuTypeDetailHead(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.MenuTypeDetailHeadReq)
	dishTypeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	menuType, err := ms.menuService.GetMenuType(req.MenuTypeID)
	if err != nil {
		res.Code = enum.SystemError
		return
	}
	head := conv.GenerateMenuTypeDetailTableHead(menuType, dishTypeMap)
	res.Data = head
}

func (ms *MenuServer) RequestMenuTypeDetailData(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.MenuTypeDetailDataReq)
	dishTypeMap, err := ms.dishService.GetDishTypeMap()
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	menuType, err := ms.menuService.GetMenuType(req.MenuTypeID)
	if err != nil {
		logger.Warn(menuServerLogTag, "RequestMenuTypeDetailData GetMenuType Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	data := conv.GenerateMenuTypeDetailTableData(menuType, dishTypeMap)
	res.Data = data
}
