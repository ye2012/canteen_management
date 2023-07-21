package service

import (
	"database/sql"

	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	storeServiceLogTag = "StoreService"
)

type StoreService struct {
	storeModel     *model.StorehouseModel
	storeTypeModel *model.StorehouseTypeModel
	goodsModel     *model.GoodsModel
	goodsTypeModel *model.GoodsTypeModel

	storeTypeMap map[uint32]*model.StorehouseType
	goodsMap     map[uint32]*model.Goods
}

func NewStoreService(sqlCli *sql.DB) *StoreService {
	storeModel := model.NewStorehouseModelWithDB(sqlCli)
	storeTypeModel := model.NewStorehouseTypeModelWithDB(sqlCli)
	goodsModel := model.NewGoodsModelWithDB(sqlCli)
	goodsTypeModel := model.NewGoodsTypeModelWithDB(sqlCli)
	return &StoreService{
		storeModel:     storeModel,
		storeTypeModel: storeTypeModel,
		goodsModel:     goodsModel,
		goodsTypeModel: goodsTypeModel,
		storeTypeMap:   make(map[uint32]*model.StorehouseType),
		goodsMap:       make(map[uint32]*model.Goods),
	}
}

func (ss *StoreService) Init() error {
	typeList, err := ss.GetStoreTypeList()
	if err != nil {
		logger.Warn(storeServiceLogTag, "Init Failed|Err:%v", err)
		return err
	}
	for _, typeInfo := range typeList {
		ss.storeTypeMap[typeInfo.ID] = typeInfo
	}

	return nil
}

func (ss *StoreService) updateStoreTypeCache(storeType *model.StorehouseType) {
	ss.storeTypeMap[storeType.ID] = storeType
}

func (ss *StoreService) GetStoreTypeList() ([]*model.StorehouseType, error) {
	typeList, err := ss.storeTypeModel.GetStorehouseTypes()
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetStoreTypeList Failed|Err:%v", err)
		return nil, err
	}
	return typeList, nil
}

func (ss *StoreService) AddStoreType(storeType *model.StorehouseType) error {
	err := ss.storeTypeModel.Insert(storeType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "Insert StoreType Failed|Err:%v", err)
		return err
	}
	ss.updateStoreTypeCache(storeType)
	return nil
}

func (ss *StoreService) UpdateStoreType(storeType *model.StorehouseType) error {
	err := ss.storeTypeModel.UpdateStorehouseType(storeType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "UpdateStoreType Failed|Err:%v", err)
		return err
	}
	ss.updateStoreTypeCache(storeType)
	return nil
}

func (ss *StoreService) GetStoreList(storeType uint32) (map[uint32][]*model.StorehouseDetail, error) {
	storeGoodsList, err := ss.storeModel.GetStorehouseGoods(storeType)
	if err != nil {
		logger.Warn(storeServiceLogTag, "GetMenuList Failed|Err:%v", err)
		return nil, err
	}

	retMap := make(map[uint32][]*model.StorehouseDetail)
	for _, storeGoods := range storeGoodsList {
		if _, ok := retMap[storeGoods.StoreTypeID]; ok == false {
			retMap[storeGoods.StoreTypeID] = make([]*model.StorehouseDetail, 0)
		}
		retMap[storeGoods.StoreTypeID] = append(retMap[storeGoods.StoreTypeID], storeGoods)
	}

	return retMap, nil
}
