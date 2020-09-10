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
	TYPE_PHONE   = 1
	// cookie存储的key
	COOKIE_NAME  = "auth"
	USER_ID      = "user_id"
)

