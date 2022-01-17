package thirdLogin

// 微信相关常量
const (
    // 微信开放平台appid
    WECHAT_APPID            = "wx9306bf43be47830b"
    // 微信开放平台secret
	WECHAT_SECRET           = "3af44d06170ecdab4b49d1c70268c71f"
	// 微信获取access token url
	WECHAT_ACCESS_TOKEN_URL = "https://api.weixin.qq.com/sns/oauth2/access_token?"
	// 微信用户信息url
	WECHAT_USER_INFO_URL    = "https://api.weixin.qq.com/sns/userinfo?"
	// 微信用户手机号
	WECHAT_USER_MOBILE_URL  = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?"
	// 小程序appid
	APPLET_APPID            = "wx668f0f98d65655ae"
	// 小程序secret
	APPLET_SECRET           = "2360683db36bb132de7b8ef0b0e8c05a"
	// 小程序获取access token url
	APPLET_ACCESS_TOKEN_URL = "https://api.weixin.qq.com/cgi-bin/token?"
	// 登录凭证校验
	APPLET_CODE2_SESSION_URL = "https://api.weixin.qq.com/sns/jscode2session?"
)

// 微博相关常量
const (
	// 微博用户信息url
	WEIBO_USER_INFO_URL     = "https://api.weibo.com/2/users/show.json?"
)

// QQ相关常量
const (
	// qq获取unionid url
	QQ_GET_UNIONID_URL      = "https://graph.qq.com/oauth2.0/me?access_token="
	// qq用户信息url
	QQ_USER_INFO_URL        = "http://openapi.tencentyun.com/v3/user/get_info?"
	// QQ iOS appkey及appid
	IOS_QQ_APP_KEY          = "DEQ6LMyBkscqe5oA"
	IOS_QQ_APP_ID           = "1106700522"
    // QQ android appkey及appid
	ANDROID_QQ_APP_KEY      = "EP9P5SCaNy1c98UR"
	ANDROID_QQ_APP_ID       = "1106668666"
	// 需区分android和iOS
	IPHONE                  = "iPhone"
	ANDROID                 = "Android"
)

