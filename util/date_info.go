package util

import (
	"time"
)

/**
 * @Description 获得当前月的初始和结束日期
 * @Param  * @param null
 * @return
 **/
func GetMonthDay(now time.Time) (string, string) {
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	f := firstOfMonth.Unix()
	l := lastOfMonth.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

/**
 * @Description 获得当前周的初始和结束日期
 * @Param  * @param null
 * @return
 **/
//func GetWeekDay(now time.Time) (string, string) {
//	offset := int(time.Monday - now.Weekday())
//	//周日做特殊判断 因为time.Monday = 0
//	if offset > 0 {
//		offset = -6
//	}
//
//	lastoffset := int(time.Saturday - now.Weekday())
//	//周日做特殊判断 因为time.Monday = 0
//	if lastoffset == 6 {
//		lastoffset = -1
//	}
//
//	firstOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
//	lastOfWeeK := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, lastoffset+1)
//	f := firstOfWeek.Unix()
//	l := lastOfWeeK.Unix()
//	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
//}

// 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
func GetWeekDay(now time.Time, week int) (startTime, endTime string) {
	
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}
	
	year, month, day := now.Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02")
	endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02")
	
	return startTime, endTime
}

/**
 * @Description //获得当前季度的初始和结束日期
 * @Param  * @param null
 * @return
 **/
func GetQuarterDay() (string, string) {
	year := time.Now().Format("2006")
	month := int(time.Now().Month())
	var firstOfQuarter string
	var lastOfQuarter string
	if month >= 1 && month <= 3 {
		//1月1号
		firstOfQuarter = year + "-01-01 00:00:00"
		lastOfQuarter = year + "-03-31 23:59:59"
	} else if month >= 4 && month <= 6 {
		firstOfQuarter = year + "-04-01 00:00:00"
		lastOfQuarter = year + "-06-30 23:59:59"
	} else if month >= 7 && month <= 9 {
		firstOfQuarter = year + "-07-01 00:00:00"
		lastOfQuarter = year + "-09-30 23:59:59"
	} else {
		firstOfQuarter = year + "-10-01 00:00:00"
		lastOfQuarter = year + "-12-31 23:59:59"
	}
	return firstOfQuarter, lastOfQuarter
}
