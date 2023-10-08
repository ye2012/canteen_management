package conv

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/model"
)

func ConvertToUserInfoList(daoList []*model.AdminUser) []*dto.UserInfo {
	retList := make([]*dto.UserInfo, 0, len(daoList))
	for _, dao := range daoList {
		retList = append(retList, &dto.UserInfo{ID: dao.ID, NickName: dao.NickName, UserName: dao.UserName,
			PhoneNumber: dao.PhoneNumber, Role: dao.Role, OpenID: dao.OpenID})
	}
	return retList
}

func ConvertFromUserInfo(info *dto.UserInfo) *model.AdminUser {
	return &model.AdminUser{ID: info.ID, NickName: info.NickName, UserName: info.UserName, Password: info.Password,
		PhoneNumber: info.PhoneNumber, Role: info.Role, OpenID: info.OpenID}
}
