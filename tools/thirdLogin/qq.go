package thirdLogin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"encoding/json"
	"net/url"
	"strings"
	"fmt"
	"errors"
)

type QQ struct {}

func NewQQ() *QQ {
	return &QQ{}
}

type QQUnionid struct {
	ClientId string `json:"client_id"`
	OpenId   string `json:"openid"`
	Unionid  string `json:"unionid"`     // 关联id
}

type QQUserInfo struct {
	Ret                int    `json:"ret"`                // 返回码
	Msg                string `json:"msg"`                // 如果ret<0，会有相应的错误信息提示，返回数据全部用UTF-8编码。
	Nickname           string `json:"nickname"`           // 用户在QQ空间的昵称。
	Gender             string `json:"gender"`             // 性别。 如果获取不到则默认返回”男”
	Country            string `json:"country"`            // 国家（当pf=qzone、pengyou或qplus时返回）。
	Province           string `json:"province"`           // 省（当pf=qzone、pengyou或qplus时返回）。
	City               string `json:"city"`               // 市（当pf=qzone、pengyou或qplus时返回）。
	Figureurl          string `json:"figureurl"`          // 大小为30×30像素的QQ空间头像URL。
	IsYellowVip        int    `json:"is_yellow_vip"`      // 标识用户是否为黄钻用户（0：不是；1：是）。
	IsYellowYearVip    int    `json:"is_yellow_year_vip"` // 标识是否为年费黄钻用户（0：不是； 1：是）
	YellowVipLevel     int    `json:"yellow_vip_level"`   // 黄钻等级
	IsYellowHighVip    int    `json:"is_yellow_high_vip"` // 是否为豪华版黄钻用户（0：不是； 1：是）。（当pf=qzone、pengyou或qplus时返回）
}

// 获取QQ关联id
func (qq *QQ) GetQQUnionID(accessToken string) (string, error) {
	resp, body, errs := gorequest.New().Get(fmt.Sprintf("%s%s&unionid=1", QQ_GET_UNIONID_URL, accessToken)).End()
	if errs != nil {
		log.Fatalf("qq_trace: get qq user info err %+v", errs)
		return "", errors.New("request error")
	}

	log.Printf("resp: %+v, body: %s, errs: %+v", resp, string(body), errs)
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("qq_trace: read all err: %v", err)
		return "", err
	}

	r := bytes.Fields(respData)
	u := new(QQUnionid)
	if err = json.Unmarshal(r[1], u); err != nil {
		log.Fatalf("qq_trace: unmarshal err: %v", err)
		return "", err
	}

	return u.Unionid, nil
}

// 获取qq app配置信息
func (qq *QQ) GetQQAppConf(agent string) (string, string) {
	var qqAppKey, qqAppId string
	switch agent {
	case IPHONE:
		qqAppKey = IOS_QQ_APP_KEY
		qqAppId = IOS_QQ_APP_ID
	case ANDROID:
		qqAppKey = ANDROID_QQ_APP_KEY
		qqAppId = ANDROID_QQ_APP_ID
	}

	return qqAppKey, qqAppId
}

// 获取QQ用户信息(agent: iPhone OR Android)
func (qq *QQ) GetQQUserInfo(openid, access_token, agent, pf string) (*QQUserInfo, error) {
	appKey, appId := qq.GetQQAppConf(agent)
	if appKey == "" || appId == "" {
		log.Fatalf("qq_trace: invalid agent, agent:%v", agent)
		return nil, errors.New("data invalid")
	}
	v := url.Values{}
	v.Set("openid", openid)
	v.Set("openkey", access_token)
	v.Set("appid", appId)
	v.Set("pf", pf)
	v.Set("format", "json")
	// 签名生成
	sig := qq.SigGenerate(v, appKey)
	v.Set("sig", sig)

	qqUserInfo := new(QQUserInfo)
	resp, body, errs := gorequest.New().Get(QQ_USER_INFO_URL + v.Encode()).EndStruct(qqUserInfo)
	if errs != nil {
		log.Fatalf("gorequest body: %+v, err: %v", string(body), errs)
		return nil, errors.New("http request failed!")
	}

	if qqUserInfo.Ret != 0 || resp.StatusCode != 200 {
		log.Fatalf("gorequest body: %+v", string(body))
		return nil, errors.New("status not 200!")
	}

	return qqUserInfo, nil
}

// QQ签名生成
func (qq *QQ) SigGenerate(v url.Values, qqAppKey string) (sig string) {
	/*Step 1. 构造源串*/
	//第1步：将请求的URI路径进行URL编码（URI不含host，URI示例：/v3/user/get_info）。
	uri := "/v3/user/get_info"
	uriUrlEncode := url.QueryEscape(uri)

	//第2步：将除“sig”外的所有参数按key进行字典升序排列。
	//第3步：将第2步中排序后的参数(key=value)用&拼接起来，并进行URL编码。
	paramUrlEncode := url.QueryEscape(v.Encode())

	//第4步：将HTTP请求方式（GET或者POST）以及第1步和第3步中的字符串用&拼接起来。
	reqMode := "GET"
	a := []string{reqMode, uriUrlEncode, paramUrlEncode}
	sourceStr := strings.Join(a, "&")

	/*Step 2. 构造密钥*/
	//得到密钥的方式：在应用的appkey末尾加上一个字节的“&”，即appkey&
	secretKey := qqAppKey + "&"

	/*Step 3. 生成签名值*/
	//1. 使用HMAC-SHA1加密算法，使用Step2中得到的密钥对Step1中得到的源串加密。
	//（注：一般程序语言中会内置HMAC-SHA1加密算法的函数，例如PHP5.1.2之后的版本可直接调用hash_hmac函数。）
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(sourceStr))

	//2. 然后将加密后的字符串经过Base64编码。
	//（注：一般程序语言中会内置Base64编码函数，例如PHP中可直接调用 base64_encode() 函数。）
	sig = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}
