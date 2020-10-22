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
	// QQ用户
	TYPE_QQ        = 3
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
	DEFAULT_AVATAR = "https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=2863307160,2081500909&fm=26&gp=0.jpg?imageView2/1/w/80/h/80"
)

const (
	// 昵称最长30个字符
	MAX_NAME_LEN      = 30
	// 昵称最少1个字符
	MIN_NAME_LEN      = 1
	// 签名最长140个字符
	MAX_SIGNATURE_LEN = 140
)

// 用户状态 0 正常 1 封禁
const (
	USER_NORMAL = 0
	USER_FORBID = 1
)

