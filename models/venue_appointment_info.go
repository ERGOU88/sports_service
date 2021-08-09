package models

type VenueAppointmentInfo struct {
	Id              int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	TimeNode        string `json:"time_node" xorm:"not null default '' comment('时间节点 例如 10:00-12:00') VARCHAR(128)"`
	Duration        int    `json:"duration" xorm:"not null default 0 comment('总时长（秒）') INT(8)"`
	RealAmount      int    `json:"real_amount" xorm:"not null comment('真实价格（单位：分）') INT(11)"`
	CurAmount       int    `json:"cur_amount" xorm:"not null comment('当前价格 (包含真实价格、 折扣价格（单位：分）') INT(11)"`
	DiscountRate    int    `json:"discount_rate" xorm:"not null default 0 comment('折扣率') INT(11)"`
	DiscountAmount  int    `json:"discount_amount" xorm:"not null default 0 comment('优惠的金额') INT(11)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('0 正常 1 不可用') TINYINT(1)"`
	QuotaNum        int    `json:"quota_num" xorm:"not null default 0 comment('配额：可容纳人数/可预约人数 -1表示没有限制') INT(10)"`
	RelatedId       int64  `json:"related_id" xorm:"not null comment('场馆id/私教课程id/大课id') BIGINT(20)"`
	RecommendType   int    `json:"recommend_type" xorm:"not null default 0 comment('推荐类型 0 无 1 热门推荐 2 低价推荐') TINYINT(1)"`
	AppointmentType int    `json:"appointment_type" xorm:"not null default 0 comment('0 场馆预约 1 私教课程预约 2 大课预约') TINYINT(1)"`
	WeekNum         int    `json:"week_num" xorm:"not null comment('1 周一 2 周二 3 周三 4 周四 5 周五 6 周六 0 周日') TINYINT(1)"`
	Sortorder       int    `json:"sortorder" xorm:"not null default 0 comment('排序权重 倒序') INT(11)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
