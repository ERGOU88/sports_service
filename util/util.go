package util

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
	"github.com/json-iterator/go"
	"github.com/rs/xid"
	"github.com/zheng-ji/goSnowFlake"
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

// NewID 年+月+日+时间+4位随机
func NewID() int64 {
	t := time.Now().Format("20060102150405")
	t = fmt.Sprintf("%s%d", t, GenerateRandnum(1000, 9999))
	id, _ := strconv.ParseInt(t, 10, 64)
	return id
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


