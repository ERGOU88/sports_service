package consts

type Duration string
// 视频时长
// 0 表示没有限制 1 表示 1～5分钟  2：5～10分钟 3：10～30分钟 4：30分钟以上
const (
	UNLIMITED_DURATION  Duration = "0"
	ONE_TO_FIVE_MINUTES Duration = "1"
	FIVE_TO_TEN_MINUTES Duration = "2"
	TEN_TO_HALF_HOUR    Duration = "3"
	MORE_THAN_HALF_HOUR Duration = "4"
)


type LimitType string
// 发布时间
// 0 不限制时间 1 一天内 2 一周内 3 半年内
const (
	UNLIMITED_TIME  LimitType = "0"
	A_DAY           LimitType = "1"
	A_WEEK          LimitType = "2"
	HALF_A_YEAR     LimitType = "3"
)

// 综合搜索默认展示3条
const (
	DEFAULT_SEARCH_PAGE = 1
	DEFAULT_SEARCH_SIZE = 3
)

