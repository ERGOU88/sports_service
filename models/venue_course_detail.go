package models

type VenueCourseDetail struct {
	Id              int64  `json:"id" xorm:"pk autoincr comment('课程id') BIGINT(20)"`
	Title           string `json:"title" xorm:"comment('课程标题') MEDIUMTEXT"`
	Subhead         string `json:"subhead" xorm:"not null default '' comment('课程副标题') VARCHAR(521)"`
	Describe        string `json:"describe" xorm:"comment('课程描述') MEDIUMTEXT"`
	PromotionPic    string `json:"promotion_pic" xorm:"not null default '' comment('宣传图') VARCHAR(1000)"`
	Icon            string `json:"icon" xorm:"not null default '' comment('图标') VARCHAR(256)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	VenueId         int64  `json:"venue_id" xorm:"not null default 0 comment('场馆id') BIGINT(20)"`
	CostDescription string `json:"cost_description" xorm:"not null default '' comment('费用说明') VARCHAR(1000)"`
	Instructions    string `json:"instructions" xorm:"not null default '' comment('购买须知') VARCHAR(1000)"`
}
