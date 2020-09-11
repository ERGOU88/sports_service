package consts

type PLATFORM int32

const (
	// android端
	ANDROID_PLATFORM    PLATFORM = iota
	// iOS端
	IOS_PLATFORM
	// web端
	WEB_PLATFORM
)

const (
	// 手机用户
	TYPE_MOBILE    = 0
	// 微信用户
	TYPE_WECHAT    = 1
	// 微博用户
	TYPE_WEIBO     = 2
)

const (
	BOY_OR_GIRL  = 0
	BOY          = 1
	GIRL         = 2
)

const (
	// cookie存储的key
	COOKIE_NAME    = "auth"
	USER_ID        = "user_id"
	// todo: 默认头像
	DEFAULT_AVATAR = ""
)

const (
	WECHAT_ACCESS_TOKEN_URL = "https://api.weixin.qq.com/sns/oauth2/access_token?"
	// 微信用户信息
	WECHAT_USER_INFO_URL    = "https://api.weixin.qq.com/sns/userinfo?"
	// 微博用户信息
	WEIBO_USER_INFO_URL     = "https://api.weibo.com/2/users/show.json?"
)

