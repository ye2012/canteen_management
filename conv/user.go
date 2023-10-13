package conv

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/model"
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
	return &model.AdminUser{ID: info.ID, NickName: info.NickName, UserName: info.UserName, Password: info.Password,
		PhoneNumber: info.PhoneNumber, Role: finalRole, OpenID: info.OpenID}
}
