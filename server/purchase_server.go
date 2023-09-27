package server

import (
	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
	"github.com/canteen_management/service"
	"github.com/canteen_management/utils"
	"github.com/gin-gonic/gin"
)

const (
	purchaseServerLogTag = "PurchaseServer"
)

type PurchaseServer struct {
	purchaseService *service.PurchaseService
}

func NewPurchaseServer(dbConf utils.Config) (*PurchaseServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(purchaseServerLogTag, "NewMysqlClient Failed|Err:%v", err)
		return nil, err
	}
	purchaseService := service.NewPurchaseService(sqlCli)
	return &PurchaseServer{
		purchaseService: purchaseService,
	}, nil
}

func (ps *PurchaseServer) RequestSupplierList(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.SupplierListReq)

	supplierList, err := ps.purchaseService.GetSupplierList(req.Name, req.PhoneNumber)
	if err != nil {
		logger.Warn(purchaseServerLogTag, "GetSupplierList Failed|Err:%v", err)
		res.Code = enum.SystemError
		return
	}

	retData := &dto.SupplierListRes{
		SupplierList: ConvertToSupplier(supplierList),
	}
	res.Data = retData
}

func (ps *PurchaseServer) RequestModifySupplier(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.ModifySupplierReq)

	switch req.Operate {
	case enum.OperateTypeAdd:
		err := ps.purchaseService.AddSupplier(ConvertFromSupplierInfo(req.Supplier))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	case enum.OperateTypeModify:
		err := ps.purchaseService.UpdateSupplier(ConvertFromSupplierInfo(req.Supplier))
		if err != nil {
			res.Code = enum.SqlError
			return
		}
	default:
		logger.Warn(purchaseServerLogTag, "RequestModifySupplier Unknown OperateType|Type:%v", req.Operate)
		res.Code = enum.SystemError
	}
}

func (ps *PurchaseServer) RequestBindSupplier(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	req := rawReq.(*dto.BindSupplierReq)

	err := ps.purchaseService.BindSupplier(req.SupplierID, req.OpenID)
	if err != nil {
		res.Code = enum.SqlError
		res.Msg = err.Error()
		return
	}
}

func ConvertToSupplier(suppliers []*model.Supplier) []*dto.SupplierInfo {
	retList := make([]*dto.SupplierInfo, 0)
	for _, supplier := range suppliers {
		retInfo := &dto.SupplierInfo{
			SupplierID:       supplier.ID,
			Name:             supplier.Name,
			PhoneNumber:      supplier.PhoneNumber,
			IDNumber:         supplier.IDNumber,
			Location:         supplier.Location,
			ValidityDeadline: supplier.ValidityDeadline.Unix(),
		}
		retList = append(retList, retInfo)
	}
	return retList
}

func ConvertFromSupplierInfo(supplierInfo *dto.SupplierInfo) *model.Supplier {
	return &model.Supplier{ID: supplierInfo.SupplierID, Name: supplierInfo.Name, PhoneNumber: supplierInfo.PhoneNumber,
		IDNumber: supplierInfo.IDNumber, Location: supplierInfo.Location}
}
