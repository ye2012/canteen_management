package conv

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/canteen_management/dto"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/model"
)

const (
	orderConvertLogTag = "OrderConvert"
)

func ConvertMenuToOrderNode(menuDate int64, dayMenu map[uint8][]uint32, dishMap map[uint32]*model.Dish,
	typeMap map[uint32]*model.DishType, dishQuantityMap map[string]float64, includeAll bool) []*dto.OrderNode {
	retData := make([]*dto.OrderNode, 0)
	for mealType := enum.MealUnknown + 1; mealType < enum.MealALL; mealType++ {
		totalDishList, ok := dayMenu[mealType]
		if !ok {
			continue
		}

		mealName := time.Unix(menuDate, 0).Format("01-02") + enum.GetMealName(mealType)
		retMeal := &dto.OrderNode{ID: fmt.Sprintf("%v_%v", menuDate, mealType), Name: mealName}
		dishListByType, maxTypeID := make(map[uint32][]*model.Dish), uint32(0)
		for _, dishID := range totalDishList {
			dishType := dishMap[dishID].DishType
			if dishType > maxTypeID {
				maxTypeID = dishType
			}
			if _, ok := dishListByType[dishType]; ok == false {
				dishListByType[dishType] = make([]*model.Dish, 0)
			}
			dishListByType[dishType] = append(dishListByType[dishType], dishMap[dishID])
		}

		retMeal.Children = make([]*dto.OrderNode, 0, len(dishListByType))
		mealSelected := int32(0)
		for dishType := uint32(1); dishType <= maxTypeID; dishType++ {
			dishList, ok := dishListByType[dishType]
			if !ok {
				continue
			}
			retListByType := &dto.OrderNode{ID: fmt.Sprintf("%v", dishType), Name: typeMap[dishType].DishTypeName}
			retListByType.Children = make([]*dto.OrderNode, 0, len(dishList))
			for index, dish := range dishList {
				retDish := &dto.OrderNode{ID: fmt.Sprintf("%v_%v_%v", retMeal.ID, dish.ID, index),
					DishID: dish.ID, Name: dish.DishName, Price: dish.Price}
				retListByType.Children = append(retListByType.Children, retDish)
				if quantity, ok := dishQuantityMap[retDish.ID]; (ok && quantity > 0) || includeAll {
					mealSelected += int32(quantity)
				}
			}
			retMeal.Children = append(retMeal.Children, retListByType)
		}
		retMeal.SelectedNumber = mealSelected
		retData = append(retData, retMeal)
	}

	return retData
}

func ConvertToOrderDao(uid uint32, phoneNumber, ID string, buildingID, floor uint32, room string) *model.OrderDao {
	ids := strings.Split(ID, "_")
	if len(ids) != 2 {
		logger.Warn(orderConvertLogTag, "ID illegal|ID:%v", ID)
		return nil
	}
	mealTime, _ := strconv.ParseInt(ids[0], 10, 32)
	mealType, _ := strconv.ParseInt(ids[1], 10, 32)

	return &model.OrderDao{
		OrderDate:   time.Unix(mealTime, 0),
		MealType:    uint8(mealType),
		Uid:         uid,
		PhoneNumber: phoneNumber,
		BuildingID:  buildingID,
		Floor:       floor,
		Room:        room,
		PayMethod:   enum.PayMethodWeChat,
	}
}

func ConvertToOrderDetailDao(items []*dto.ApplyItem) []*model.OrderDetail {
	retList := make([]*model.OrderDetail, 0)
	for _, item := range items {
		retInfo := &model.OrderDetail{
			DishID:   item.DishID,
			Quantity: item.Quantity,
		}
		retList = append(retList, retInfo)
	}
	return retList
}

func ConvertToOrderInfoList(orderList []*model.OrderDao, detailMap map[uint32][]*model.OrderDetail,
	dishMap map[uint32]*model.Dish) []*dto.OrderInfo {
	retList := make([]*dto.OrderInfo, 0)
	for _, order := range orderList {
		retInfo := &dto.OrderInfo{
			ID:            fmt.Sprintf("%v_%v", order.OrderDate.Unix(), order.MealType),
			Name:          order.OrderDate.Format("01-02") + enum.GetMealName(order.MealType),
			OrderID:       fmt.Sprintf("%v", order.ID),
			OrderNo:       "",
			PayOrderID:    order.PayOrderID,
			MealType:      order.MealType,
			UserPhone:     order.PhoneNumber,
			BuildingID:    order.BuildingID,
			Floor:         order.Floor,
			Room:          order.Room,
			TotalAmount:   order.TotalAmount,
			PayMethod:     order.PayMethod,
			PaymentAmount: order.PayAmount,
			OrderItems:    make([]*dto.ApplyItem, 0),
			CreateTime:    order.CreateAt.Unix(),
			OrderStatus:   order.Status,
		}
		if details, ok := detailMap[order.ID]; ok {
			orderItems := make([]*dto.ApplyItem, 0, len(details))
			for _, detail := range details {
				item := &dto.ApplyItem{
					DishID:   detail.DishID,
					DishName: dishMap[detail.DishID].DishName,
					Price:    detail.Price,
					Quantity: detail.Quantity,
				}
				orderItems = append(orderItems, item)
			}
			retInfo.OrderItems = orderItems
		}
		retList = append(retList, retInfo)
	}
	return retList
}

func ConvertDishID(itemID string) (uint32, error) {
	ids := strings.Split(itemID, IndexDelimiter)
	if len(ids) != 4 {
		return 0, fmt.Errorf("id不合法")
	}

	goodsID, err := strconv.ParseInt(ids[2], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("转化菜品id失败|ID:%v", itemID)
	}
	return uint32(goodsID), nil
}
