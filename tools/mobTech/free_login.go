package mobTech

import (
	"encoding/json"
	"log"
	"time"
	"errors"
)

// 移动开发者服务平台
type MobTech struct {
	*base
}

// 返回的结构
type Data struct {
	Status int         `json:"status"`
	Error  string      `json:"error"`
	Res    interface{} `json:"res"`
}

// mob一键登陆返回值
type FreeLoginRes struct {
	IsValid        int32      `json:"isValid"`        // 验证状态，1:成功, 2:失败
	Phone          string     `json:"phone"`          // 返回手机号
}

// 实栗
func NewMobTech() *MobTech {
	return &MobTech{
		&base{},
	}
}

// 一键登陆 服务端校验
// appkey app标识
// token 客户端的token
// opToken 客户端返回的运营商token
// operator 客户端返回的运营商，CMCC:中国移动通信, CUCC:中国联通通讯, CTCC:中国电信
// sign 签名（MD5(所有参数使用key的正序，通过a=b&b=c+appSecret组成)）
// timestamp 当前时间戳（毫秒）
// 返回的数据需要使用appSecet解密
// DES/CBC/PKCS5Padding 偏移量 00000000 使用base64转码
func (mob *MobTech) FreeLogin(token, opToken, operator string) (string, error) {
	data := map[string]interface{}{
		"appkey":    APP_KEY,
		"token":     token,
		"opToken":   opToken,
		"operator":  operator,
		"timestamp": time.Now().UnixNano() / 1e6,
	}

	data["sign"] = mob.generateSign(data, APP_SECRET)
	postBody, err := mob.HttpPostBody(FREE_LOGIN_URL, data)
	if err != nil {
		log.Fatalf("mob_trace: http request failed, err:%v", err)
		return "", err
	}

	ret := new(Data)
	if err := json.Unmarshal(postBody, &ret); err != nil {
		log.Fatalf("mob_trace: unmarshal err:%v", err)
		return "", err
	}

	if ret.Status != 200 {
		log.Fatal("mob_trace: request status not 200")
		return "", errors.New("request status not 200")
	}

	res := new(FreeLoginRes)
	decode, _ := mob.Base64Decode([]byte(ret.Res.(string)))
	decr, _ := mob.DesDecrypt(decode, []byte(APP_SECRET)[0:8])
	if err := json.Unmarshal(decr, &res); err != nil {
		return "", err
	}

	if res.IsValid != SUCCESS {
		return "", errors.New("free login failed")
	}

	return res.Phone, nil
}


