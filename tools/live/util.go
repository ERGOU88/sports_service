package live

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

// 生成txTime
func BuildTxTime(expireTm int64) string {
	timestamp := time.Now().Unix()
	timestamp = timestamp + expireTm
	return strconv.FormatInt(timestamp, 16)
}

// 生成txSecret
func BuildTxSecret(key, roomId string, expireTm int64) string {
	streamName := roomId
	txTime := BuildTxTime(expireTm)
	secretStr := key + streamName + txTime
	return Md5(secretStr)
}

func Md5(str string) string {
	data := []byte(str)
	hash := md5.New()
	hash.Write(data)
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr
}
