package utils

import "time"

func GetFirstDateOfWeek(curTime int64) int64 {
	cur := time.Unix(curTime, 0)
	offset := int(time.Monday - cur.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(cur.Year(), cur.Month(), cur.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStartDate.Unix()
}

func GetZeroTime(curTime int64) int64 {
	cur := time.Unix(curTime, 0)
	return time.Date(cur.Year(), cur.Month(), cur.Day(), 0, 0, 0, 0, time.Local).Unix()
}

func GetMidDayTime(curTime int64) int64 {
	cur := time.Unix(curTime, 0)
	return time.Date(cur.Year(), cur.Month(), cur.Day(), 12, 0, 0, 0, time.Local).Unix()
}
