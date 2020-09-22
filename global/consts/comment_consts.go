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

// 评论内容最少10个字符 最多1000个字符
const (
	COMMENT_MIN_LEN  = 10
	COMMENT_MAX_LEN  = 1000
)
