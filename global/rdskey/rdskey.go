package rdskey

import (
	"fmt"
)

const (
	KEY_EXPIRE_MIN    = 60 * 1
	KEY_EXPIRE_HOUR   = 60 * 60
	KEY_EXPIRE_DAY    = 60 * 60 * 24
	KEY_EXPIRE_WEEK   = KEY_EXPIRE_DAY * 7
)

const (
	SUGAR = "fpv:"
)

const (
	LOGIN_REPEAT       = SUGAR + "login_repeat_type:%d_sid:%s"              // 拦截重复注册的问题{拼接设备类型 + 手机号码/unionid}
	USER_AUTH          = SUGAR + "key_user_auth_%s"                         // 保存用户token {拼接user_id}
	USER_NICKNAME_INCR = SUGAR + "user_nickname_incr"                       // 用户昵称自增
)

// make redis key
func MakeKey(key_fmt string, keys ...interface{}) string {
	return fmt.Sprintf(key_fmt, keys...)
}
