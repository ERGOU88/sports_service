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
	VIDEO_INVALID_LABEL_NAME    = 10011
	VIDEO_ADD_VIDEO_LABEL_FAIL  = 10012
	VIDEO_LABEL_NOT_EXISTS      = 10013
	VIDEO_LABEL_DELETE_FAIL     = 10014
	VIDEO_BARRAGE_DELETE_FAIL   = 10015

	// 11001 - 12000 评论相关错误码
	COMMENT_NOT_EXISTS          = 11001
	COMMENT_DELETE_FAIL         = 11002

	// 12001 - 13000 配置相关错误码
	CONFIG_ADD_BANNER_FAIL      = 12001
	CONFIG_DEL_BANNER_FAIL      = 12002
	CONFIG_ADD_AVATAR_FAIL      = 12003
	CONFIG_DEL_AVATAR_FAIL      = 12004

	// 13001 - 14000 用户相关错误码
	USER_FORBID_FAIL            = 13001
	USER_UNFORBID_FAIL          = 13002
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
	VIDEO_INVALID_LABEL_NAME:    "视频标签名称长度不合法，应小于等于20个字符",
	VIDEO_ADD_VIDEO_LABEL_FAIL:  "添加视频标签失败",
	VIDEO_LABEL_NOT_EXISTS:      "视频标签不存在",
	VIDEO_LABEL_DELETE_FAIL:     "视频标签删除失败",
  VIDEO_BARRAGE_DELETE_FAIL:   "视频弹幕删除失败",

	COMMENT_NOT_EXISTS:          "评论不存在",
	COMMENT_DELETE_FAIL:         "删除评论失败",

	CONFIG_ADD_BANNER_FAIL:      "添加banner失败",
	CONFIG_DEL_BANNER_FAIL:      "删除banner失败",
	CONFIG_ADD_AVATAR_FAIL:      "添加系统头像失败",
	CONFIG_DEL_AVATAR_FAIL:      "删除系统头像失败",

	USER_FORBID_FAIL:            "用户封禁失败",
	USER_UNFORBID_FAIL:          "用户解封失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}



