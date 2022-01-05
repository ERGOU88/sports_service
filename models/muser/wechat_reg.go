package muser

type wechatRegister struct {
	*base
}

// 微信登陆 请求参数
type WxLoginParam struct {
	Code    string     `binding:"required" json:"code" example:"code码"`
}

// 小程序登录 请求参数
type AppletLoginParam struct {
	Code          string   `binding:"required" json:"code"`            // 用户登录凭证（有效期五分钟）
	//CodeByPhone   string   `binding:"required" json:"code_by_phone"`   // 手机号获取凭证
	PhoneData     string   `binding:"required" json:"phone_data"`
	Iv            string   `binding:"required" json:"iv"`
}

type BindWechatParam struct {
	Code          string   `binding:"required" json:"code"`            // 用户登录凭证（有效期五分钟）
	UserId        string
}


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

// 实栗
func NewWechatRegister() *wechatRegister {
	return &wechatRegister{
		&base{},
	}
}




