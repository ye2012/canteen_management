package server

import (
	"database/sql"
	"github.com/canteen_management/enum"
	"time"

	"github.com/canteen_management/dto"
	"github.com/canteen_management/logger"
	"github.com/canteen_management/service"
	"github.com/canteen_management/utils"
	"github.com/gin-gonic/gin"
)

const (
	statisticLogTag = "StatisticServer"
)

type StatisticServer struct {
	sqlCli       *sql.DB
	orderService *service.OrderService
}

func NewStatisticServer(dbConf utils.Config) (*StatisticServer, error) {
	sqlCli, err := utils.NewMysqlClient(dbConf)
	if err != nil {
		logger.Warn(orderServerLogTag, "NewOrderServer Failed|Err:%v", err)
		return nil, err
	}
	orderService := service.NewOrderService(sqlCli)
	return &StatisticServer{
		sqlCli:       sqlCli,
		orderService: orderService,
	}, nil
}

func (ss *StatisticServer) RequestDashboard(ctx *gin.Context, rawReq interface{}, res *dto.Response) {
	resData := &dto.DashboardRes{}
	curTime, err := time.Now().Unix(), error(nil)
	err = ss.GetMonthOrderInfo(curTime, resData)
	if err != nil {
		logger.Warn(statisticLogTag, "RequestDashboard GetMonthOrderInfo Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	err = ss.GetTodayOrderInfo(curTime, resData)
	if err != nil {
		logger.Warn(statisticLogTag, "RequestDashboard GetTodayOrderInfo Failed|Err:%v", err)
		res.Code = enum.SqlError
		return
	}

	res.Data = resData
}

func (ss *StatisticServer) GetTodayOrderInfo(curTime int64, resData *dto.DashboardRes) error {
	start, end := utils.GetDayTimeRange(curTime)
	sqlStr := "SELECT SUM(pay_amount), COUNT(*), `status`, `meal_type` from `order` WHERE `order_date` >= ? AND `order_date` <= ? GROUP BY `meal_type`, `status`"
	rows, err := ss.sqlCli.Query(sqlStr, time.Unix(start, 0), time.Unix(end, 0))
	if err != nil {
		logger.Warn(statisticLogTag, "GetMonthOrderInfo Query Failed|Err:%v", err)
		return err
	}
	defer rows.Close()
	payAmount, count, status, mealType := 0.0, int32(0), 0, uint8(0)

	for rows.Next() {
		err = rows.Scan(&payAmount, &count, &status, &mealType)
		if err != nil {
			logger.Warn(statisticLogTag, "RequestDashboard Scan Failed|Err:%v", err)
			continue
		}
		resData.DayOrderCount += count
		if status != enum.OrderNew && status != enum.OrderCancel {
			switch mealType {
			case enum.MealBreakfast:
				resData.DayBreakfastCount += count
			case enum.MealLunch:
				resData.DayLunchCount += count
			case enum.MealDinner:
				resData.DayDinnerCount += count
			}
			resData.DaySuccessCount += count
			resData.DayPayAmount += payAmount
		}

	}
	return nil
}

func (ss *StatisticServer) GetMonthOrderInfo(curTime int64, resData *dto.DashboardRes) error {
	start := utils.GetMonthStartTime(curTime)
	end := utils.GetDayEndTime(curTime)
	sqlStr := "SELECT SUM(pay_amount), COUNT(*), `status` from `order` WHERE `order_date` >= ? AND `order_date` <= ?  GROUP BY `status`"
	rows, err := ss.sqlCli.Query(sqlStr, time.Unix(start, 0), time.Unix(end, 0))
	if err != nil {
		logger.Warn(statisticLogTag, "GetMonthOrderInfo Query Failed|Err:%v", err)
		return err
	}
	defer rows.Close()
	payAmount, count, status := 0.0, int32(0), 0

	for rows.Next() {
		err = rows.Scan(&payAmount, &count, &status)
		if err != nil {
			logger.Warn(statisticLogTag, "RequestDashboard Scan Failed|Err:%v", err)
			continue
		}
		resData.TotalOrderCount += count
		if status != enum.OrderNew && status != enum.OrderCancel {
			resData.TotalSuccessCount += count
			resData.TotalPayAmount += payAmount
		}
	}
	return nil
}
