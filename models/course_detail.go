package models

type CourseDetail struct {
	Id               int64  `json:"id" xorm:"pk autoincr comment('课程id') BIGINT(20)"`
	Title            string `json:"title" xorm:"comment('课程标题') MEDIUMTEXT"`
	Describe         string `json:"describe" xorm:"comment('课程描述') MEDIUMTEXT"`
	SaiCoin          int    `json:"sai_coin" xorm:"not null default 0 comment('课程价格（代币数 * 100存储）存在小数') INT(11)"`
	EventSaiCoin     int    `json:"event_sai_coin" xorm:"not null default 0 comment('活动价格 (代币数）') INT(11)"`
	EventStartTime   int    `json:"event_start_time" xorm:"not null default 0 comment('活动开始时间') INT(11)"`
	EventEndTime     int    `json:"event_end_time" xorm:"not null default 0 comment('活动结束时间') INT(11)"`
	PromotionPic     string `json:"promotion_pic" xorm:"not null default '' comment('宣传图') VARCHAR(521)"`
	TeacherPhoto     string `json:"teacher_photo" xorm:"not null default '' comment('老师照片') VARCHAR(256)"`
	TeacherTitle     string `json:"teacher_title" xorm:"not null default '' comment('老师title') VARCHAR(512)"`
	TeacherName      string `json:"teacher_name" xorm:"not null default '' comment('老师名称') VARCHAR(64)"`
	Icon             string `json:"icon" xorm:"not null default '' comment('图标地址') VARCHAR(256)"`
	Sortorder        int    `json:"sortorder" xorm:"not null default 1 comment('排序') INT(11)"`
	IsRecommend      int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop            int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	IsFree           int    `json:"is_free" xorm:"not null default 0 comment('0 收费 1 免费') TINYINT(1)"`
	VipIsFree        int    `json:"vip_is_free" xorm:"not null default 0 comment('会员是否免费 0 收费 1 免费') TINYINT(1)"`
	CreateAt         int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt         int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Status           int    `json:"status" xorm:"not null default 0 comment('0 正常 1 待发布 2 剔除') TINYINT(1)"`
	AreasOfExpertise string `json:"areas_of_expertise" xorm:"not null default '' comment('老师擅长领域') VARCHAR(512)"`
}
