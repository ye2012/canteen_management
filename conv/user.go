package conv

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/model"
	"github.com/canteen_management/utils"
)

func ConvertToRoleList(role uint32) []uint32 {
	roleList := make([]uint32, 0)
	for i := enum.RoleMin + 1; i < enum.RoleMax; i++ {
		if role&(1<<i) > 0 {
			roleList = append(roleList, uint32(i))
		}
	}
	return roleList
}

func ConvertToUserInfoList(daoList []*model.AdminUser) []*dto.UserInfo {
	retList := make([]*dto.UserInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.UserInfo{ID: dao.ID, NickName: dao.NickName, UserName: dao.UserName,
			PhoneNumber: dao.PhoneNumber, RoleList: ConvertToRoleList(dao.Role), OpenID: dao.OpenID})
	}
	return retList
}

func ConvertFromUserInfo(info *dto.UserInfo) *model.AdminUser {
	finalRole := uint32(0)
	for _, role := range info.RoleList {
		finalRole += 1 << role
	}
	hashedPass := utils.Encrypt(info.Password)
	return &model.AdminUser{ID: info.ID, NickName: info.NickName, UserName: info.UserName, Password: hashedPass,
		PhoneNumber: info.PhoneNumber, Role: finalRole, OpenID: info.OpenID}
}

func ConvertToRouterTypeInfoList(daoList []*model.RouterType) []*dto.RouterTypeInfo {
	retList := make([]*dto.RouterTypeInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.RouterTypeInfo{RouterTypeID: dao.ID, RouterTypeName: dao.RouterTypeName,
			SortID: dao.SortID})
	}
	return retList
}

func ConvertFromRouterTypeInfo(info *dto.RouterTypeInfo) *model.RouterType {
	return &model.RouterType{ID: info.RouterTypeID, RouterTypeName: info.RouterTypeName, SortID: info.SortID}
}

func ConvertToRouterInfoList(daoList []*model.RouterDetail) []*dto.RouterInfo {
	retList := make([]*dto.RouterInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.RouterInfo{RouterID: dao.ID, RouterType: dao.RouterType, RouterName: dao.RouterName,
			RouterPath: dao.RouterPath, RouterSortID: dao.RouterSortID, RoleList: ConvertToRoleList(dao.Role)})
	}
	return retList
}

func ConvertFromRouterInfo(info *dto.RouterInfo) *model.RouterDetail {
	finalRole := uint32(2)
	for _, role := range info.RoleList {
		finalRole |= 1 << role
	}
	return &model.RouterDetail{ID: info.RouterID, RouterType: info.RouterType, RouterName: info.RouterName,
		RouterPath: info.RouterPath, RouterSortID: info.RouterSortID, Role: finalRole}
}

func ConvertToRouterNode(routerList []*model.RouterDetail, routerTypeList []*model.RouterType, role uint32) []*dto.RouterNode {
	retList := make([]*dto.RouterNode, 0, len(routerTypeList))
	routerMap := make(map[uint32][]*model.RouterDetail)
	for _, router := range routerList {
		if router.Role&role == 0 {
			continue
		}
		if _, ok := routerMap[router.RouterType]; ok == false {
			routerMap[router.RouterType] = make([]*model.RouterDetail, 0)
		}
		routerMap[router.RouterType] = append(routerMap[router.RouterType], router)
	}
	for _, routerType := range routerTypeList {
		routers, ok := routerMap[routerType.ID]
		if ok == false || len(routers) == 0 {
			continue
		}
		retInfo := &dto.RouterNode{Name: routerType.RouterTypeName}
		for _, router := range routers {
			routerInfo := &dto.RouterNode{Name: router.RouterName, Path: router.RouterPath}
			retInfo.Children = append(retInfo.Children, routerInfo)
		}
		retList = append(retList, retInfo)
	}
	return retList
}
