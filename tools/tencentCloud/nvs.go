package tencentCloud

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"errors"
	"fmt"
	"sports_service/server/util"
	"time"
)

const (
	TENCENT_SDK_APP_ID  = 1400521069
	TENCENT_NVS_URL     = "https://yun.tim.qq.com/v5/rapidauth/validate?sdkappid=%d&random=%d"
)

// 返回的结构
type Response struct {
	Result    int         `json:"result"`   // 0表示成功，非0表示失败
 	Errmsg    string      `json:"errmsg"`   // 错误信息
	Mobile    string      `json:"mobile"`   // 手机号码
}

// 一键登录 校验客户端token 获取手机号码
// carrier 运营商，移动：mobile， 联通：unicom，电信：telecom
// nationcode 86表示中国大陆
func (tc *TencentCloud) FreeLogin(token, carrier, nationcode string) (string, error) {
	random := util.GetXID()
	tm := time.Now().Unix()
	data := map[string]interface{}{
		"sdkappid":  TENCENT_SDK_APP_ID,
		"sig":       tc.generateSign(tm, random),
		"carrier":   carrier,
		"token":     token,
		"random":    random,
		"nationcode": nationcode,
		"time":  tm,
	}

	postBody, err := tc.HttpPostBody(fmt.Sprintf(TENCENT_NVS_URL, TENCENT_SDK_APP_ID, random), data)
	if err != nil {
		log.Printf("tencent_trace: http request failed, err:%v", err)
		return "", err
	}

	log.Printf("postBody:%s", string(postBody))

	ret := new(Response)
	if err := util.JsonFast.Unmarshal(postBody, &ret); err != nil {
		log.Printf("tencent_trace: unmarshal err:%v", err)
		return "", err
	}

	return ret.Mobile, nil

}

// 生成签名
func (tc *TencentCloud) generateSign(tm, random int64) string {
	ret := fmt.Sprintf("appkey=%s&random=%d&time=%d", tc.secretKey, random, tm)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(ret))
	cipherStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

// post请求
func (tc *TencentCloud) HttpPostBody(url string, msg map[string]interface{}) ([]byte, error) {
	request := gorequest.New()
	resp, body, errs := request.Post(url).Set("Content-Type", "application/json; charset=utf-8").SendMap(msg).End()
	log.Printf("resp:%+v", resp)
	if errs != nil {
		log.Printf("mob_trace: request err:%+v, body:%s", errs, body)
		return nil, errors.New("request error")
	}

	return ioutil.ReadAll(resp.Body)
}
