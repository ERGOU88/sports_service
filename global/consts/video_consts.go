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
