package consts

// 状态 -1 查看所有 0 审核中 1 审核成功 2 审核失败 3 逻辑删除
const (
	POST_VIEW_ALL      = "-1"
	POST_UNDER_REVIEW  = "0"
	POST_AUDIT_SUCCESS = "1"
	POST_AUDIT_FAILURE = "2"
	POST_DELETE_STATUS = "3"
)

// 社区帖子来源类型
// 0 社区发布 1 转发视频  2 转发帖子
const (
	COMMUNITY_PUB_POST      = iota
	COMMUNITY_FORWARD_VIDEO
	COMMUNITY_FORWARD_POST
)

// 社区帖子内容类型
// 0 纯文本 1 图文 2 视频+文本
const (
	POST_TYPE_TEXT  = iota
	POST_TYPE_IMAGE
	POST_TYPE_VIDEO
)


// 排序方式 0 时间 1 热度 （点赞 + 浏览 + 评论）
const (
	POST_SORT_TIME   = "0"
	POST_SORT_HOT    = "1"
)

