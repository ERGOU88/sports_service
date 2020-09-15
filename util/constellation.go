package util

import (
	"fmt"
	"time"
)

// 解析出年、月、日
func GetTimeFromStrDate(date string) (year, month, day int) {
	const shortForm = "2006-01-02"
	d, err := time.Parse(shortForm, date)
	if err != nil {
		fmt.Println("出生日期解析错误！")
		return 0, 0, 0
	}

	year = d.Year()
	month = int(d.Month())
	day = d.Day()
	return
}

// 获取年龄
func GetAge(year, month, day int) (age int) {
	if year <= 0 || month <= 0 || month >= 13 || day <= 0 || day >= 32 {
		age = -1
	}

	nowyear := time.Now().Year()
	nmonth := time.Now().Month()
	nowmonth := int(nmonth)
	nowday := time.Now().Day()
	if nowmonth > month {
		age = nowyear - year
	} else if nowmonth < month {
		age = nowyear - year - 1
	} else if nowmonth == month {
		if nowday >= day {
			age = nowyear - year
		} else {
			age = nowyear - year - 1
		}
	}

	return
}
