package errdef

const (
	SUCCESS            = 200
	ERROR              = 500
	INVALID_PARAMS     = 400

	// 10001 - 11000 视频相关错误码
	VIDEO_ALREADY_DELETE        = 10001
	VIDEO_NOT_EXISTS            = 10002
	VIDEO_ALREADY_PASS          = 10003
	VIDEO_EDIT_STATUS_FAIL      = 10004
	VIDEO_DELETE_PUBLISH_FAIL   = 10005
	VIDEO_DELETE_LABEL_FAIL     = 10006
	VIDEO_DELETE_STATISTIC_FAIL = 10007
	VIDEO_NOT_PASS              = 10008
	VIDEO_EDIT_TOP_FAIL         = 10009
	VIDEO_EDIT_RECOMMEND_FAIL   = 10010
	VIDEO_INVALID_LABEL_LEN     = 10011
  VIDEO_INVALID_LABEL_NAME    = 10012
	VIDEO_ADD_VIDEO_LABEL_FAIL  = 10013
	VIDEO_LABEL_NOT_EXISTS      = 10014
	VIDEO_LABEL_DELETE_FAIL     = 10015
	VIDEO_BARRAGE_DELETE_FAIL   = 10016

	// 11001 - 12000 评论相关错误码
	COMMENT_NOT_EXISTS          = 11001
	COMMENT_DELETE_FAIL         = 11002

	// 12001 - 13000 配置相关错误码
	CONFIG_ADD_BANNER_FAIL      = 12001
	CONFIG_INVALID_END_TIME     = 12002
	CONFIG_DEL_BANNER_FAIL      = 12003
	CONFIG_ADD_AVATAR_FAIL      = 12004
	CONFIG_DEL_AVATAR_FAIL      = 12005
	CONFIG_ADD_HOT_SEARCH_FAIL  = 12006
  CONFIG_DEL_HOT_SEARCH_FAIL  = 12007
  CONFIG_SET_STATUS_HOT_FAIL  = 12008
  CONFIG_SET_SORT_HOT_FAIL    = 12009
  CONFIG_INVALID_HOT_SEARCH   = 12010
  CONFIG_HOT_NAME_EXISTS      = 12011
  CONFIG_COS_ACCESS_FAIL      = 12012
  CONFIG_ADD_PACKAGE_FAIL     = 12013
  CONFIG_UPDATE_PACKAGE_FAIL  = 12014
  CONFIG_DEL_PACKAGE_FAIL     = 12015
  CONFIG_UPDATE_BANNER_FAIL   = 12016

	// 13001 - 14000 用户相关错误码
	USER_FORBID_FAIL            = 13001
	USER_UNFORBID_FAIL          = 13002

	// 14001 - 15000 后台用户相关错误码
	ADMIN_ADD_FAIL              = 14001
	ADMIN_HAS_EXISTS            = 14002
	ADMIN_NOT_EXISTS            = 14003
	ADMIN_PASSWORD_NOT_MATCH    = 14004

  // 15001 - 16000 通知相关错误码
  NOTIFY_INVALID_SEND_TM      = 15001
  NOTIFY_INVALID_CONTENT      = 15002
  NOTIFY_INVALID_USER_IDS     = 15003
  NOTIFY_USER_NOT_FOUND       = 15004
  NOTIFY_PUSH_FAIL            = 15005
  NOTIFY_MSG_NOT_EXISTS       = 15006
  NOTIFY_CAN_NOT_CANCEL       = 15007
  NOTIFY_CANCEL_FAIL          = 15008
  NOTIFY_INVALID_START_TM     = 15009
  NOTIFY_CAN_NOT_DEL          = 15010
  NOTIFY_DEL_FAIL             = 15011

)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	VIDEO_ALREADY_DELETE:        "视频已删除",
	VIDEO_NOT_EXISTS:            "视频不存在",
	VIDEO_ALREADY_PASS:          "视频已过审，只能执行删除操作",
	VIDEO_EDIT_STATUS_FAIL:      "视频状态更新失败",
	VIDEO_DELETE_PUBLISH_FAIL:   "删除发布的视频失败",
	VIDEO_DELETE_LABEL_FAIL:     "删除视频标签失败",
	VIDEO_DELETE_STATISTIC_FAIL: "删除视频统计数据失败",
	VIDEO_NOT_PASS:              "视频未审核",
	VIDEO_EDIT_TOP_FAIL:         "修改视频置顶状态失败",
	VIDEO_EDIT_RECOMMEND_FAIL:   "修改视频推荐状态失败",
  VIDEO_INVALID_LABEL_LEN:     "视频标签名称长度不合法，应小于等于20个字符",
	VIDEO_INVALID_LABEL_NAME:    "标签名称中含有违规字符",
	VIDEO_ADD_VIDEO_LABEL_FAIL:  "添加视频标签失败",
	VIDEO_LABEL_NOT_EXISTS:      "视频标签不存在",
	VIDEO_LABEL_DELETE_FAIL:     "视频标签删除失败",
  VIDEO_BARRAGE_DELETE_FAIL:   "视频弹幕删除失败",

	COMMENT_NOT_EXISTS:          "评论不存在",
	COMMENT_DELETE_FAIL:         "删除评论失败",

	CONFIG_ADD_BANNER_FAIL:      "添加banner失败",
  CONFIG_INVALID_END_TIME:     "上架时间必须大于下架时间，下架时间必须大于当前时间（且上架时长需>=30分钟）",
	CONFIG_DEL_BANNER_FAIL:      "删除banner失败",
	CONFIG_ADD_AVATAR_FAIL:      "添加系统头像失败",
	CONFIG_DEL_AVATAR_FAIL:      "删除系统头像失败",
  CONFIG_ADD_HOT_SEARCH_FAIL:  "添加热搜失败",
  CONFIG_DEL_HOT_SEARCH_FAIL:  "删除热搜失败",
  CONFIG_SET_STATUS_HOT_FAIL:  "设置热搜状态失败",
  CONFIG_SET_SORT_HOT_FAIL:    "设置热搜权重失败",
  CONFIG_INVALID_HOT_SEARCH:   "热搜词含有违规文字",
  CONFIG_HOT_NAME_EXISTS:      "热搜词已存在",
  CONFIG_COS_ACCESS_FAIL:      "获取cos通行证失败",
  CONFIG_ADD_PACKAGE_FAIL:     "添加新包失败",
  CONFIG_UPDATE_PACKAGE_FAIL:  "更新包信息失败",
  CONFIG_DEL_PACKAGE_FAIL:     "删除包信息失败",
  CONFIG_UPDATE_BANNER_FAIL:   "banner更新失败",

  USER_FORBID_FAIL:            "用户封禁失败",
	USER_UNFORBID_FAIL:          "用户解封失败",

  ADMIN_ADD_FAIL:              "管理员添加失败",
  ADMIN_HAS_EXISTS:            "管理员已存在",
  ADMIN_NOT_EXISTS:            "管理员不存在",
  ADMIN_PASSWORD_NOT_MATCH:    "帐号/密码不正确",

  NOTIFY_INVALID_SEND_TM:         "无效的发送时间，发送时间必须大于当前时间",
  NOTIFY_INVALID_CONTENT:         "推送内容含有违规文字",
  NOTIFY_INVALID_USER_IDS:        "无效的用户id列表",
  NOTIFY_USER_NOT_FOUND:          "用户未查找到",
  NOTIFY_PUSH_FAIL:               "系统消息推送失败",
  NOTIFY_MSG_NOT_EXISTS:          "系统通知不存在",
  NOTIFY_CAN_NOT_CANCEL:          "该通知无法撤回",
  NOTIFY_CANCEL_FAIL:             "系统通知撤回失败",
  NOTIFY_INVALID_START_TM:        "最大可设置当前时间 + 7天",
  NOTIFY_CAN_NOT_DEL:             "无法删除！请先撤回该通知",
  NOTIFY_DEL_FAIL:                "删除失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}



