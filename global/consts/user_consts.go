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
	USER_STATUS    = "user_status"
	// todo: 默认头像
	DEFAULT_AVATAR = "http://qafpv-backend.youzu.com/upload/2020_11_02/543195997677817856.png"
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

// 后台用户列表排序 0 按注册时间倒序 1 关注数 2 粉丝数 3 发布数 4 浏览数 5 点赞数 6 收藏数 7 评论数 8 弹幕数
const (
  USER_SORT_BY_TIME      = "0"
  USER_SORT_BY_ATTENTION = "1"
  USER_SORT_BY_FANS      = "2"
  USER_SORT_BY_PUBLISH   = "3"
  USER_SORT_BY_BROWSE    = "4"
  USER_SORT_BY_LIKE      = "5"
  USER_SORT_BY_COLLECT   = "6"
  USER_SORT_BY_COMMENT   = "7"
  USER_SORT_BY_BARRAGE   = "8"
)

