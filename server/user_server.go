package server

import (
	"github.com/canteen_management/config"
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
	userServerLogTag = "UserServer"
)

type UserServer struct {
	userService  *service.UserService
	orderService *service.OrderService
}

func NewUserServer(dbConf utils.Config) (*UserServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	userService := service.NewUserService(sqlCli)
	orderService := service.NewOrderService(sqlCli)
	if err != nil {
		logger.Warn(userServerLogTag, "NewUserServer Failed|Err:%v", err)
		return nil, err
	}

	return &UserServer{
		userService:  userService,
		orderService: orderService,
	}, nil
}

func (us *UserServer) RequestBindPhoneNumber(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.BindPhoneNumberReq)

	err := us.userService.BindPhoneNumber(req.Uid, req.PhoneNumber)
	if err != nil {
		res.Code = enum.SystemError
		res.Msg = err.Error()
		return
	}
	return
}

func (us *UserServer) RequestCanteenLogin(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.CanteenLoginReq)

	openID, err := utils.MiniProgramLogin(config.Config.AppID, config.Config.AppSecret, req.Code)
	if err != nil {
		res.Code = enum.SystemError
		res.Msg = err.Error()
		return
	}

	user, err := us.userService.WxUserLogin(openID)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
	role := us.userService.GetWxUserRole(openID)
	resData := &dto.CanteenLoginRes{
		Uid:         user.ID,
		OpenID:      user.OpenID,
		PhoneNumber: user.PhoneNumber,
		Role:        role,
	}
	discountType := us.userService.GetWxUserDiscount(user.OpenID)
	resData.ExtraPay, resData.Discount, resData.DiscountLeft, err = us.orderService.LoginUserOrderDiscountInfo(user.ID, discountType)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}

	res.Data = resData
}

func (us *UserServer) RequestKitchenLoginLogin(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.KitchenLoginReq)

	openID, err := utils.MiniProgramLogin(config.Config.AppID, config.Config.AppSecret, req.Code)
	if err != nil {
		res.Code = enum.SystemError
		res.Msg = err.Error()
		return
	}

	user, err := us.userService.WxUserLogin(openID)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
	role := us.userService.GetWxUserRole(openID)
	resData := &dto.KitchenLoginRes{
		Uid:         user.ID,
		OpenID:      user.OpenID,
		PhoneNumber: user.PhoneNumber,
		Role:        role,
	}

	res.Data = resData
}

func (us *UserServer) RequestOrderUserList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.OrderUserListReq)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize > 1000 {
		req.PageSize = 100
	}
	userList, userNumber, err := us.userService.GetOrderUserList(req.OpenID, req.PhoneNumber, req.DiscountLevel, req.Page, req.PageSize)
	if err != nil {
		logger.Warn(userServerLogTag, "GetOrderUserList Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	userInfoList := make([]*dto.OrderUserInfo, 0, len(userList))
	for _, user := range userList {
		userInfo := &dto.OrderUserInfo{
			ID:            user.ID,
			OpenID:        user.OpenID,
			PhoneNumber:   user.PhoneNumber,
			DiscountLevel: user.DiscountLevel,
		}
		userInfoList = append(userInfoList, userInfo)
	}

	retData := &dto.OrderUserListRes{
		UserList: userInfoList,
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: int32(userNumber),
		},
	}
	res.Data = retData
}

func (us *UserServer) RequestModifyOrderUser(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
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
		err := us.userService.AddOrderUser(userList)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := us.userService.UpdateOrderUser(userList[0])
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(userServerLogTag, "RequestModifyOrderUser Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (us *UserServer) RequestBindOrderUser(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.BindOrderUserReq)

	err := us.userService.BindOrderUser(req.ID, req.OpenID)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func (us *UserServer) RequestAdminUserList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.AdminUserListReq)
	list, count, err := us.userService.GetAdminUserList(req.RoleType, req.Page, req.PageSize)
	if err != nil {
		res.Code = enum.SqlError
		return
	}

	res.Data = &dto.UserListRes{
		UserList: conv.ConvertToUserInfoList(list),
		PaginationRes: dto.PaginationRes{
			Page:        req.Page,
			PageSize:    req.PageSize,
			TotalNumber: count,
		},
	}
}

func (us *UserServer) RequestModifyAdminUser(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifyAdminUserReq)
	userInfo := conv.ConvertFromUserInfo(req.User)
	switch req.Operate {
	case enum.OperateTypeAdd:
		err := us.userService.AddAdminUser(userInfo)
		if err != nil {
			res.Code = enum.SqlError
			res.Msg = err.Error()
			return
		}
	case enum.OperateTypeModify:
		err := us.userService.UpdateAdminUser(userInfo)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeDel:
		err := us.userService.DeleteAdminUser(userInfo.ID)
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(userServerLogTag, "RequestModifyUser Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (us *UserServer) RequestBindAdminUser(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.BindAdminReq)
	err := us.userService.BindAdminUser(req.User.ID, req.User.OpenID)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}
