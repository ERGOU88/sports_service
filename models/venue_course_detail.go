package models

type VenueCourseDetail struct {
	Id             int64  `json:"id" xorm:"pk autoincr comment('课程id') BIGINT(20)"`
	CoachId        int64  `json:"coach_id" xorm:"not null default 0 comment('关联的教兽id [私教课才会关联]') index BIGINT(20)"`
	ClassPeriod    int    `json:"class_period" xorm:"not null comment('单课程时长（秒）') INT(11)"`
	Title          string `json:"title" xorm:"comment('课程标题') MEDIUMTEXT"`
	Subhead        string `json:"subhead" xorm:"not null default '' comment('课程副标题') VARCHAR(521)"`
	Describe       string `json:"describe" xorm:"comment('课程描述') MEDIUMTEXT"`
	Price          int    `json:"price" xorm:"not null comment('课程价格（分/每课时）') INT(11)"`
	EventPrice     int    `json:"event_price" xorm:"not null default 0 comment('活动价格') INT(11)"`
	EventStartTime int    `json:"event_start_time" xorm:"not null default 0 comment('活动开始时间') INT(11)"`
	EventEndTime   int    `json:"event_end_time" xorm:"not null default 0 comment('活动结束时间') INT(11)"`
	PromotionPic   string `json:"promotion_pic" xorm:"not null default '' comment('宣传图') VARCHAR(1000)"`
	Icon           string `json:"icon" xorm:"not null default '' comment('图标') VARCHAR(256)"`
	Sortorder      int    `json:"sortorder" xorm:"not null default 1 comment('排序') INT(11)"`
	IsRecommend    int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop          int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Status         int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	CourseType     int    `json:"course_type" xorm:"not null default 0 comment('1 私教课 2 大课') TINYINT(1)"`
	PeriodNum      int    `json:"period_num" xorm:"not null comment('总课时数') INT(6)"`
	VenueId        int64  `json:"venue_id" xorm:"not null default 0 comment('场馆id') BIGINT(20)"`
}
