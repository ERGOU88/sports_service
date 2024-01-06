package im

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"net/url"
	"sports_service/util"
)

const (
	// 根据运行环境 走配置
	TX_IM_APP_ID  = 1400576334
	TX_IM_APP_KEY = "080789c06a28b355e1ec94b97ad61edc4eb887275e68fa0046f5659cead396e7"
	// 控制台配置的管理员
	TX_IDENTIFIER = "bluetrans"
	TX_IM_HOST    = "https://console.tim.qq.com"
)

// 生成签名
//func GenSig(userId string, expireTm int) (string, error) {
//	sig, err := tencentyun.GenUserSig(config.Global.TencentImAppId, config.Global.TencentImSecret, userId, expireTm)
//	if err != nil {
//		return "", err
//	}
//
//	return sig, nil
//}

// 生成请求url
func GenRequestUrl(appId int, identifier, sig, uri string) string {
	values := url.Values{}
	values.Add("sdkappid", fmt.Sprint(appId))
	values.Add("identifier", identifier)
	values.Add("random", fmt.Sprint(util.GenerateRandnum(100000, 999999)))
	values.Add("usersig", sig)
	values.Add("contenttype", "json")

	return fmt.Sprintf("%s%s?%s", TX_IM_HOST, uri, values.Encode())
}

// post请求
func HttpPostBody(url string, msg map[string]interface{}) ([]byte, error) {
	request := gorequest.New()
	resp, body, errs := request.Post(url).Set("Content-Type", "application/json; charset=utf-8").SendMap(msg).End()
	log.Printf("resp:%+v", resp)
	if errs != nil {
		log.Printf("im_trace: request err:%+v, body:%s", errs, body)
		return nil, errors.New("request error")
	}

	return ioutil.ReadAll(resp.Body)
}
