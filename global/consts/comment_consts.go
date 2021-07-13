package consts

// 发布评论为1级评论  回复评论为2级评论
const (
	COMMENT_PUBLISH    = 1
	COMMENT_REPLY      = 2
)

// 排序方式 0 时间 1 热度
const (
	SORT_TIME   = "0"
	SORT_HOT    = "1"
)

// 评论内容最少1个字符 最多1000个字符
const (
	COMMENT_MIN_LEN  = 1
	COMMENT_MAX_LEN  = 1000
)

// 评论列表排序 1 按评论时间倒序 2 按评论点赞数 3 按评论回复数
const (
	SORT_BY_TIME  = "1"
	SORT_BY_LIKE  = "2"
	SORT_BY_REPLY = "3"
)

// 评论类型 0 视频 1 帖子
const (
	COMMENT_TYPE_VIDEO  = iota
	COMMENT_TYPE_POST
)
