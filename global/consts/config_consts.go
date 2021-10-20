package consts

type WorkType int

// 1 视频 2 帖子 3 资讯
const (
	WORK_TYPE_VIDEO WorkType = 1
	WORK_TYPE_POST  WorkType = 2
	WORK_TYPE_INFO  WorkType = 3
)


// 用户行为类型
// 1 点赞 2 浏览 3 分享 4 弹幕 5 评论 6 收藏
const (
	ACTION_TYPE_FABULOUS = 1
	ACTION_TYPE_BROWSE   = 2
	ACTION_TYPE_SHARE    = 3
	ACTION_TYPE_BARRAGE  = 4
	ACTION_TYPE_COMMENT  = 5
	ACTION_TYPE_COLLECT  = 6
)
