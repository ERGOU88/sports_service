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
	USER_ID_INCR              = SUGAR + "user_id_incr"                             // 用户id自增
	USER_READ_BELIKED_NOTIFY  = SUGAR + "user_read_beliked_notify_%s"              // 记录用户读取被点赞消息的最新时间{拼接user_id}
	USER_READ_AT_NOTIFY       = SUGAR + "user_read_at_notify_%s"                   // 记录用户读取被@消息的最新时间{拼接user_id}
	USER_READ_ATTENTION_VIDEO = SUGAR + "user_read_attention_pub_%s"               // 记录用户读取关注用户发布的视频的最新时间（刷新列表时才记录）{拼接user_id}

	SMS_INTERVAL_NUM          = SUGAR + "sms:interval_num:%s_%s"                   // 一天内同一手机发送验证码次数{拼接年月日_手机号码}
	SMS_INTERVAL_TM           = SUGAR + "sms:interval_tm:%s_%s"                    // 验证码间隔时间60秒 {拼接短信类型_手机号}
	SMS_CODE                  = SUGAR + "sms:code:%s_%s"                           // 验证码内容{拼接短信类型_手机号}

	USER_WATCHING_VIDEO       = SUGAR + "user_watching_video_%s"                   // 记录正在观看视频的用户标示[xid]{拼接视频id}

	VIDEO_UPLOAD_TASK         = SUGAR + "video_upload_task_%d"                     // 记录任务id 对应的 用户id {拼接任务id（唯一）}
	VIDEO_UPLOAD_INFO         = SUGAR + "video_upload_info_%s_%d"                  // 记录用户上传的视频信息{拼接 userId + taskId}

	SEARCH_HISTORY_CONTENT    = SUGAR + "search_history_content_%s"                // 记录历史搜索内容{拼接 userId}

	MSG_EVENT_KEY             = SUGAR + "msg_event_key"                            // 消息（App推送等）

	AUDIT_MODE                = SUGAR + "audit_mode"                               // [视频、帖子] 审核模式 1 人工 + AI 2 人工审核
)

// make redis key
func MakeKey(key_fmt string, keys ...interface{}) string {
	return fmt.Sprintf(key_fmt, keys...)
}
