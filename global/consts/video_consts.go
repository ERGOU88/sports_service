package consts

// 0为浏览的视频记录
const (
	TYPE_BROWSE_VIDEOS  = 0
)

// 0为后台管理员发布视频  1为用户发布视频
const (
	PUBLISH_VIDEO_BY_MANAGER = 0
	PUBLISH_VIDEO_BY_USER    = 1
)

// 确认操作 +1  取消操作 -1 (点赞、收藏等)
const (
    CONFIRM_OPERATE = 1
    CANCEL_OPERATE  = -1
)

// 状态查询发布的内容
// 状态 -1 查看所有 0 查看审核中 1 查看审核成功 2 查看审核失败
const (
	VIDEO_VIEW_ALL      = "-1"
	VIDEO_UNDER_REVIEW  = "0"
	VIDEO_AUDIT_SUCCESS = "1"
	VIDEO_AUDIT_FAILURE = "2"
)


// 条件查询发布的内容
// -1 发布时间 0 播放数 1 弹幕数 2 评论数 3 点赞数 4 分享数
const (
	VIDEO_CONDITION_TIME    = "-1"
	VIDEO_CONDITION_PLAY    = "0"
	VIDEO_CONDITION_BARRAGE = "1"
	VIDEO_CONDITION_COMMENT = "2"
	VIDEO_CONDITION_LIKE    = "3"
	VIDEO_CONDITION_SHARE   = "4"
)


const (
	// 默认发布时间排序
	CONDITION_FIELD_TIME    = "create_at"
	// 浏览数（播放数）排序
	CONDITION_FIELD_PLAY    = "browse_num"
	// 弹幕数排序
	CONDITION_FIELD_BARRAGE = "barrage_num"
	// 评论数排序
	CONDITION_FIELD_COMMENT = "comment_num"
	// 点赞数排序
	CONDITION_FIELD_LIKE    = "fabulous_num"
	// 分享数排序
	CONDITION_FIELD_SHARE   = "share_num"
)
