package thirdLogin

import (
	"log"
	"net/url"
	"github.com/parnurzeal/gorequest"
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
		return nil
	}

	if accessToken.Unionid == "" {
		log.Printf("err body: %s, resp: %+v", string(body), resp)
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


