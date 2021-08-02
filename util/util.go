package util

import (
	"crypto/md5"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/rs/xid"
	"github.com/zheng-ji/goSnowFlake"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"
	"strings"
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

// NewUserId 年+月+日+时间+4位随机
func NewUserId() string {
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


