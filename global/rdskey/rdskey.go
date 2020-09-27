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
	LOGIN_REPEAT              = SUGAR + "login_repeat_type:%d_sid:%s"              // 拦截重复注册的问题{拼接设备类型 + 手机号码/unionid}
	USER_AUTH                 = SUGAR + "key_user_auth_%s"                         // 保存用户token {拼接user_id}
	USER_NICKNAME_INCR        = SUGAR + "user_nickname_incr"                       // 用户昵称自增
	USER_READ_BELIKED_NOTIFY  = SUGAR + "user_read_beliked_notify_%s"              // 记录用户读取被点赞消息的最新时间{拼接user_id}
	USER_READ_AT_NOTIFY       = SUGAR + "user_read_at_notify_%s"                   // 记录用户读取被@消息的最新时间{拼接user_id}

	SMS_INTERVAL_NUM          = SUGAR + "sms:interval_num:%s_%s"                   // 一天内同一手机发送验证码次数{拼接年月日_手机号码}
	SMS_INTERVAL_TM           = SUGAR + "sms:interval_tm:%s_%s"                    // 验证码间隔时间60秒 {拼接短信类型_手机号}
	SMS_CODE                  = SUGAR + "sms:code:%s_%s"                           // 验证码内容{拼接短信类型_手机号}
)

// make redis key
func MakeKey(key_fmt string, keys ...interface{}) string {
	return fmt.Sprintf(key_fmt, keys...)
}
