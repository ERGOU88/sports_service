package thirdLogin

import (
	"encoding/json"
	"log"
	"net/url"
	"github.com/parnurzeal/gorequest"
	"errors"
	"fmt"
	log2 "sports_service/server/global/app/log"
)

type Wechat struct {}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

type WechatUserInfo struct {
	Errcode    int      `json:"errcode"`    //
	Errmsg     string   `json:"errmsg"`     //
	Openid     string   `json:"openid"`     // 普通用户的标识，对当前开发者帐号唯一
	Nickname   string   `json:"nickname"`   // 普通用户昵称
	Sex        int      `json:"sex"`        // 普通用户性别，1为男性，2为女性
	Province   string   `json:"province"`   // 普通用户个人资料填写的省份
	City       string   `json:"city"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country"`    // 国家，如中国为CN
	Headimgurl string   `json:"headimgurl"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
	Privilege  []string `json:"privilege"`  // 用户特权信息，json数组
	Unionid    string   `json:"unionid"`    // 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
}

type AppletAccessToken struct {
	Errcode      int      `json:"errcode"`    //
	Errmsg       string   `json:"errmsg"`     //
	ExpiresIn    int64    `json:"expires_in"`
	AccessToken  string   `json:"access_token"`
}

// 微信用户手机信息
type WechatPhoneInfo struct {
	Errcode      int      `json:"errcode"`    //
	Errmsg       string   `json:"errmsg"`     //
	PhoneInfo    struct {
		PhoneNumber       string   `json:"phoneNumber"`       // 用户绑定的手机号（国外手机号会有区号）
		PurePhoneNumber   string   `json:"purePhoneNumber"`   // 没有区号的手机号
		CountryCode       string   `json:"countryCode"`       // 区号
	} `json:"phone_info"`
}

type AppletCode2SessionResp struct {
	Errcode      int      `json:"errcode"`    //
	Errmsg       string   `json:"errmsg"`     //
	SessionKey   string   `json:"session_key"`         // 会话密钥
	Openid       string   `json:"openid"`              // 用户唯一标识
	Unionid      string   `json:"unionid"`             // 用户在开放平台的唯一标识符
}

type AppletUserInfo struct {
	NickName  string                 `json:"nickName"`
	OpenId    string                 `json:"openId"`
	Img       string                 `json:"avatarUrl"`
	UnionId   string                 `json:"unionId"`
	Sex       int64                  `json:"sex"`
	City      string                 `json:"city"`
	Province  string                 `json:"province"`
	Country   string                 `json:"country"`
	Watermark map[string]interface{} `json:"watermark"`
}

// 微信实栗
func NewWechat() *Wechat {
	return &Wechat{}
}

// 获取微信 access token (code 客户端获取的授权码)
func (wx *Wechat) GetWechatAccessToken(code string) *AccessToken {
	v := url.Values{}
	v.Set("code", code)
	// 开放平台appid
	v.Set("appid", WECHAT_APPID)
	// 开放平台secret
	v.Set("secret", WECHAT_SECRET)
	v.Set("grant_type", "authorization_code")
	// 返回值
	accessToken := AccessToken{}
	resp, body, errs := gorequest.New().Get(WECHAT_ACCESS_TOKEN_URL + v.Encode()).EndStruct(&accessToken)
	if errs != nil {
		log.Printf("%+v", errs)
		log2.Log.Errorf("get access token failed, err:%+v", errs)
		return nil
	}

	if accessToken.Unionid == "" {
		log.Printf("err body: %s, resp: %+v", string(body), resp)
		log2.Log.Error("unionId empty")
		return nil
	}

	return &accessToken
}

// 获取微信用户信息
func (wx *Wechat) GetWechatUserInfo(accessToken *AccessToken) *WechatUserInfo {
	v := url.Values{}
	v.Set("access_token", accessToken.AccessToken)
	v.Set("openid", accessToken.Openid)
	wxinfo := WechatUserInfo{}
	resp, body, errs := gorequest.New().Get(WECHAT_USER_INFO_URL + v.Encode()).EndStruct(&wxinfo)
	if errs != nil {
		log.Printf("get wxinfo err %+v", errs)
		return nil
	}

	log.Println("\nwxUserInfo: ", wxinfo)
	log.Println("\nresp: ", resp)
	log.Println("\nbody: ", string(body))

	if wxinfo.Errcode != 0 || resp.StatusCode != 200 {
		log.Printf("wx_trace: request failed, errCode:%d, statusCode:%d", wxinfo.Errcode, resp.StatusCode)
		return nil
	}

	return &wxinfo
}

// 获取小程序全局唯一后台接口调用凭据
func (wx *Wechat) GetAppletAccessToken() *AppletAccessToken {
	v := url.Values{}
	v.Set("grant_type", "client_credential")
	v.Set("appid", APPLET_APPID)
	v.Set("secret", APPLET_SECRET)
	info := AppletAccessToken{}
	resp, body, errs := gorequest.New().Get(APPLET_ACCESS_TOKEN_URL + v.Encode()).EndStruct(&info)
	if errs != nil {
		log.Printf("get wxinfo err %+v", errs)
		return nil
	}
	
	log.Println("\ninfo: ", info)
	log.Println("\nresp: ", resp)
	log.Println("\nbody: ", string(body))
	
	if info.Errcode != 0 || resp.StatusCode != 200 {
		log.Printf("wx_trace: request failed, errCode:%d, statusCode:%d", info.Errcode, resp.StatusCode)
		return nil
	}
	
	return &info
}

// 获取用户手机号
func (wx *Wechat) GetUserPhoneNumber(code, accessToken string) (string, error) {
	v := url.Values{}
	v.Set("access_token", accessToken)
	mp := map[string]interface{}{"code": code}
	
	info := WechatPhoneInfo{}
	resp, body, errs := gorequest.New().Post(WECHAT_USER_MOBILE_URL + v.Encode()).
		Set("Content-Type", "application/json; charset=utf-8").SendMap(mp).EndStruct(&info)
	if errs != nil {
		log.Printf("get wxinfo err %+v", errs)
		return "", errs[0]
	}
	
	log.Println("\ninfo: ", info)
	log.Println("\nresp: ", resp)
	log.Println("\nbody: ", string(body))
	
	if info.Errcode != 0 || resp.StatusCode != 200 {
		log.Printf("wx_trace: request failed, errCode:%d, statusCode:%d", info.Errcode, resp.StatusCode)
		return "", errors.New("request failed")
	}
	
	return info.PhoneInfo.PurePhoneNumber, nil
}

// 小程序登录凭证校验
func (wx *Wechat) AppletCode2Session(code string) (*AppletCode2SessionResp, error) {
	v := url.Values{}
	v.Set("js_code", code)
	// 开放平台appid
	v.Set("appid", APPLET_APPID)
	// 开放平台secret
	v.Set("secret", APPLET_SECRET)
	v.Set("grant_type", "authorization_code")
	// 返回值
	info := AppletCode2SessionResp{}
	resp, body, errs := gorequest.New().Get(APPLET_CODE2_SESSION_URL + v.Encode()).EndStruct(&info)
	if errs != nil {
		log.Printf("%+v", errs)
		return nil, errs[0]
	}
	
	log.Println("\ninfo: ", info)
	log.Println("\nresp: ", resp)
	log.Println("\nbody: ", string(body))
	
	if info.Errcode != 0 || resp.StatusCode != 200 {
		log.Printf("wx_trace: request failed, errCode:%d, statusCode:%d", info.Errcode, resp.StatusCode)
		return nil, errors.New("request failed")
	}
	
	return &info, nil
}

func (wx *Wechat) DecryptAppletUserInfo(encryptData, sessionKey, iv string) (*AppletUserInfo, error) {
	decrypt, err := AesDecrypt(encryptData, sessionKey, iv)
	if err != nil {
		return nil, err
	}
	
	userInfo := &AppletUserInfo{}
	err = json.Unmarshal(decrypt, userInfo)
	if err != nil {
		return nil, err
	}
	
	if userInfo.Watermark == nil {
		err = errors.New(fmt.Sprintf("wx login err:%s", "Watermark wrong nil"))
		return nil, err
	}
	
	appId, ok := userInfo.Watermark["appid"]
	if !ok {
		err = errors.New(fmt.Sprintf("wx login err:%s", "Watermark wrong app id not found"))
		return nil, err
	}
	
	appId = fmt.Sprintf("%v", appId)
	if appId != APPLET_APPID {
		err = errors.New(fmt.Sprintf("wx login err:%s", "Watermark wrong app id not match"))
		return nil, err
	}
	
	
	return userInfo, nil
}
