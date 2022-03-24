package consts

// 0为用户未学习某课程 1为用户学习过某课程 2 用户某课程已学习完毕
const (
	NOT_STUDY_COURSE     = 0
	HAS_STUDY_COURSE     = 1
	END_STUDY_COURSE     = 2
)

// 课程分类
// 10 表示限时免费课程
const (
	LIMITED_FREE_COURSE  = "10"
)

// 用户是否已购买课程
// 0 未购买
// 1 已购买
const (
	COURSE_NOT_PURCHASED  = 0
	COURSE_HAS_PURCHASED  = 1
)

// 课程是否免费(面向所有用户)
// 0 不免费
// 1 免费
const (
	COURSE_NOT_FREE       = 0
	COURSE_IS_FREE        = 1
)

// todo 暂时没有vip用户相关逻辑
// 课程是否免费(面向vip用户)
// 0 vip用户不免费
// 1 vip用户免费
const (
	COURSE_NOT_FREE_BY_VIP = 0
	COURSE_IS_FREE_BY_VIP  = 1
)

// 课程某一视频是否免费（面向所有用户）
// 0 不免费
// 1 免费
const (
	COURSE_VIDEO_NOT_FREE  = 0
	COURSE_VIDEO_IS_FREE   = 1
)

// todo 暂时没有vip用户相关逻辑
// 课程某一个视频是否免费(面向vip用户)
// 0 vip用户不免费
// 1 vip用户免费
const (
	COURSE_VIDEO_NOT_FREE_BY_VIP = 0
	COURSE_VIDEO_IS_FREE_BY_VIP  = 1
)

// 课程是否开启活动
// 0 没有活动
// 1 已开启活动
const (
	COURSE_NO_ACTIVITY    = 0
	COURSE_HAS_ACTIVITY   = 1
)
