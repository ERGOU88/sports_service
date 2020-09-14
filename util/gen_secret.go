package util

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NUM_STR  = "0123456789"
	CHAR_STR = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	SPEC_STR = "+=-@#~,.[]()!%^*$"
)

const (
	NUM_MODE     = "num"
	CHAR_MODE    = "char"
	MIX_MODE     = "mix"
	ADVANCE_MODE = "advance"
)

// num:只使用数字[0-9],
// char:只使用英文字母[a-zA-Z],
// mix:使用数字和字母，
// advance:使用数字、字母以及特殊字符`
func generateSecret(charset string, length int) string {
	// 初始化密码切片
	var secret = make([]byte, length, length)
	// 源字符串
	var sourceStr string
	switch charset {
	case NUM_MODE:
		sourceStr = NUM_STR
	case CHAR_MODE:
		sourceStr = CHAR_STR
	case MIX_MODE:
		sourceStr = fmt.Sprintf("%s%s", NUM_STR, CHAR_STR)
	case ADVANCE_MODE:
		sourceStr = fmt.Sprintf("%s%s%s", NUM_STR, CHAR_STR, SPEC_STR)
	default:
		sourceStr = CHAR_STR
	}

	fmt.Println("source:", sourceStr)

	// 遍历，生成一个随机index索引,
	for i := 0; i < length; i++ {
		index := rand.Intn(len(sourceStr))
		secret[i] = sourceStr[index]
	}

	return string(secret)
}

// 生成密钥
func GenSecret(mode string, length int) string {
	// 随机种子
	rand.Seed(time.Now().UnixNano())
	return generateSecret(mode, length)
}

// 生成appid 长度为8
func GenAppId() string {
	rand.Seed(time.Now().UnixNano())
	return generateSecret(MIX_MODE, 8)
}

