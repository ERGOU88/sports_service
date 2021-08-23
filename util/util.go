package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/rs/xid"
	"github.com/zheng-ji/goSnowFlake"
	"log"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	iw *goSnowFlake.IdWorker
)

var JsonFast = jsoniter.ConfigCompatibleWithStandardLibrary

// InitSnowId 初始化iw
func InitSnowId() {
	var err error
	iw, err = goSnowFlake.NewIdWorker(1)
	if err != nil {
		log.Printf("init iw err:%s", err.Error())
		return
	}
}

// GetSnowId 获得一个唯一id
func GetSnowId() int64 {
	id, _ := iw.NextId()
	return id
}

// NewOrderId 年+月+日+时间+4位随机
func NewOrderId() string {
	t := time.Now().Format("060102150405")
	return fmt.Sprintf("%s%d", t, GenerateRandnum(1000, 9999))
}

func GenerateRandnum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(max - min)
	num = num + min
	return num
}

func GetTransactionId() string {
	return fmt.Sprintf("%s%d", "FPV", GetSnowId())
}

func GetXID() int64 {
	xidService := xid.New()
	return int64(xidService.Counter())
}

func MD5(str string) string {
	md := md5.New()
	md.Write([]byte(str))
	return fmt.Sprintf("%x", md.Sum(nil))
}

func Md5String(s string) (md5_str string) {
	md5_str = fmt.Sprintf("%x", md5.Sum([]byte(s)))
	return
}

func Contains(str string, s []string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}


func PageInfo(page, size string) (p, s int) {
	p, _ = strconv.Atoi(page)
	s, _ = strconv.Atoi(size)

	if p <= 0 {
		p = 1
	}

	if s <= 0 || s > 50 {
		s = 10
	}

	return
}

func GetStrLen(r []rune) int {
	if len(r) == 0 {
		return 0
	}

	var letterlen, wordlen int
	for _, v := range r {
		// 是否为汉字
		if unicode.Is(unicode.Han, v) {
			wordlen++
			continue
		}

		letterlen++
	}

	length := letterlen + wordlen * 2
	return length
}

func TrimHtml(src string) string {
	// 将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	// 去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	// 去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	// 去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	// 去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	src = strings.Replace(src, "\n", "", -1)
	return strings.TrimSpace(src)
}

// 转化为中文展示 小于100不展示
func TransferChinese(num int) string{
	count := num/100
	if count <= 0 {
		return "0"
	}

	chineseMap:= []string{"百", "千", "万", "十万", "百万", "千万", "亿", "十亿", "百亿", "千亿"}
	chinese := fmt.Sprintf("%d%s", count/int(Pow(float64(10), len(fmt.Sprint(count))-1)), chineseMap[len(fmt.Sprint(count))-1])

	return chinese
}

func Pow(x float64, n int) float64 {
	if x == 0 {
		return 0
	}
	result := calPow(x, n)
	if n < 0 {
		result = 1 / result
	}
	return result
}

func calPow(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}

	// 向右移动一位
	result := calPow(x, n>>1)
	result *= result

	// 如果n是奇数
	if n&1 == 1 {
		result *= x
	}

	return result
}

// 检查map里面是否存在某个key，返回bool
func MapExist(m map[string]interface{}, key string) bool {
	if _, ok := m[key]; ok {
		return true
	}

	return false
}

// 解析json字符串成 map
func JsonStringToMap(jsonStr string) (m map[string]interface{}, err error) {
	mp := map[string]interface{}{}
	if err := JsonFast.Unmarshal([]byte(jsonStr), &mp); err != nil {
		return nil, err
	}

	return mp, nil
}

type DateInfo struct {
	Date      string    `json:"date"`
	Week      int       `json:"week"`
	Id        int       `json:"id"`
	WeekCn    string    `json:"week_cn"`
}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
func GetBetweenDates(start, end string) []DateInfo {
	var list []DateInfo
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(start) {
		timeFormatTpl = timeFormatTpl[0:len(start)]
	}

	date, err := time.Parse(timeFormatTpl, start)
	if err != nil {
		// 时间解析，异常
		return list
	}

	date2, err := time.Parse(timeFormatTpl, end)
	if err != nil {
		// 时间解析，异常
		return list
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return list
	}

	var id int = 1
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	info := DateInfo{}
	info.Date = date.Format("01-02")
	info.Week = int(date.Weekday())
	info.WeekCn = GetWeekCn(info.Week)
	info.Id = id
	list = append(list, info)

	for {
		id++
		info := DateInfo{}
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		info.Date = date.Format("01-02")
		info.Week = int(date.Weekday())
		info.WeekCn = GetWeekCn(info.Week)
		info.Id = id
		list = append(list, info)
		if dateStr == date2Str {
			break
		}
	}

	return list
}

// 获取周几（中文）
func GetWeekCn(week int) string {
	switch week {
	case 0:
		return "周日"
	case 1:
		return "周一"
	case 2:
		return "周二"
	case 3:
		return "周三"
	case 4:
		return "周四"
	case 5:
		return "周五"
	case 6:
		return "周六"
	}

	return ""
}

// 数字是否相邻（允许乱序）
func IsContinuous(slice []int64, n int) (bool, int64, int64) {
	if len(slice) == 0 || n == 0 {
		return false, 0, 0
	}

	max := slice[0]
	min := slice[0]

	for i := 1; i < n; i ++ {
		if slice[i] == 0 {
			return false, 0, 0
		}

		if min > slice[i] {
			min = slice[i]
		}

		if max < slice[i] {
			max = slice[i]
		}

	}

	if max - min > int64(n) - 1 {
		return false, min, max
	}

	return true, min, max
}

// 隐藏手机号码（4位）
func HideMobileNum(mobileNum string) string {
	slice := strings.Split(mobileNum, "")
	if len(slice) != 11 {
		return mobileNum
	}

	return strings.Join(slice[0:3], "") + "****" + strings.Join(slice[7:], "")
}

// 判断字符串是否都为空格
func IsSpace(r []rune) bool {
	if len(r) == 0 {
		return false
	}

	for _, v := range r {
		if b := unicode.IsSpace(v); !b {
			return true
		}
	}

	return false
}

// 转换为时分秒
func ResolveTime(seconds int) string {
	oneSecond := 60
	oneHour := 60 * oneSecond
	oneDay := 24 * oneHour

	var day = seconds / oneDay
	hour := (seconds - day * oneDay) / oneHour
	minute := (seconds - day * oneDay - hour * oneHour) / oneSecond
	second := seconds - day * oneDay - hour * oneHour - minute * oneSecond

	var res string
	if hour > 0 {
		res = fmt.Sprintf("%dh", hour)
	}

	if minute > 0 {
		res += fmt.Sprintf("%dm", minute)
	}

	if second > 0 {
		res += fmt.Sprintf("%ds", second)
	}

	return res
}

func FormatDuration(duration time.Time) string {
	hours := time.Now().Sub(duration).Hours()
	if hours > 24 * 365 {
		return duration.Format("2006-01-02")
	}

	if hours > 24 {
		return duration.Format("01-02 ")
	}

	if hours > 1 {
		return fmt.Sprintf("%d小时前", int(math.Ceil(hours)))
	}

	minute := time.Now().Sub(duration).Minutes()
	return fmt.Sprintf("%d分钟之前", int(math.Ceil(minute)))
}

// map to struct
func ToStruct(m map[string]interface{}, u interface{}) error {
	v := reflect.ValueOf(u)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Kind() != reflect.Struct {
			return errors.New("must struct")
		}

		findFromMap := func(key string,nameTag string ) interface {} {
			for k,v := range m {
				if k == key || k == nameTag {
					return v
				}
			}

			return nil
		}

		for i := 0; i < v.NumField(); i++ {
			val := findFromMap(v.Type().Field(i).Name,v.Type().Field(i).Tag.Get("name"))
			if val != nil && reflect.ValueOf(val).Kind() == v.Field(i).Kind() {
				v.Field(i).Set(reflect.ValueOf(val))
			}
		}
	}

	return errors.New("must ptr")
}
