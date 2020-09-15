package util

import (
	"crypto/md5"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/rs/xid"
	"github.com/zheng-ji/goSnowFlake"
	"log"
	"math/rand"
	"time"
	"strings"
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
		log.Fatalf("init iw err:%s", err.Error())
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
	t := time.Now().Format("20060102150405")
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


