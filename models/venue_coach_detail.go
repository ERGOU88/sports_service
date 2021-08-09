package models

type VenueCoachDetail struct {
	Id               int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	Title            string `json:"title" xorm:"comment('抬头') MEDIUMTEXT"`
	Name             string `json:"name" xorm:"not null default '' comment('私教名称') VARCHAR(60)"`
	Address          string `json:"address" xorm:"not null comment('私教地点') VARCHAR(128)"`
	Designation      string `json:"designation" xorm:"not null default '' comment('认证称号') VARCHAR(60)"`
	Describe         string `json:"describe" xorm:"comment('描述') MEDIUMTEXT"`
	AreasOfExpertise string `json:"areas_of_expertise" xorm:"not null default '' comment('擅长领域') VARCHAR(512)"`
	Cover            string `json:"cover" xorm:"not null default '' comment('封面 ') VARCHAR(256)"`
	Avatar           string `json:"avatar" xorm:"not null default '' comment('头像') VARCHAR(256)"`
	Price            int    `json:"price" xorm:"not null default 0 comment('私教价格（分）') INT(11)"`
	EventPrice       int    `json:"event_price" xorm:"not null default 0 comment('活动价格 (分)') INT(11)"`
	EventStartTime   int    `json:"event_start_time" xorm:"not null default 0 comment('活动开始时间') INT(11)"`
	EventEndTime     int    `json:"event_end_time" xorm:"not null default 0 comment('活动结束时间') INT(11)"`
	Sortorder        int    `json:"sortorder" xorm:"not null default 1 comment('排序权重') INT(11)"`
	IsRecommend      int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop            int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	CreateAt         int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt         int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Status           int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	CoachType        int    `json:"coach_type" xorm:"not null default 0 comment('1 私教课老师 2 大课老师') TINYINT(1)"`
	CourseId         int64  `json:"course_id" xorm:"not null default 0 comment('课程id') BIGINT(20)"`
}
