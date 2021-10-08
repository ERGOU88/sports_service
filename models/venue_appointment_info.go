package models

import (
	"time"
)

type VenueAppointmentInfo struct {
	Id              int64     `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	TimeNode        string    `json:"time_node" xorm:"not null default '' comment('时间节点 例如 10:00-12:00') VARCHAR(128)"`
	Duration        int       `json:"duration" xorm:"not null default 0 comment('总时长（秒）') INT(8)"`
	RealAmount      int       `json:"real_amount" xorm:"not null comment('真实价格（单位：分）') INT(11)"`
	CurAmount       int       `json:"cur_amount" xorm:"not null comment('当前价格 (包含真实价格、 折扣价格（单位：分）') INT(11)"`
	DiscountRate    int       `json:"discount_rate" xorm:"not null default 0 comment('折扣率 例如9.5折 则值为95') INT(11)"`
	DiscountAmount  int       `json:"discount_amount" xorm:"not null default 0 comment('优惠的金额') INT(11)"`
	Status          int       `json:"status" xorm:"not null default 0 comment('0 正常 1 不可用') TINYINT(4)"`
	QuotaNum        int       `json:"quota_num" xorm:"not null default 0 comment('配额：可容纳人数/可预约人数 -1表示没有限制') INT(10)"`
	RecommendTag    int       `json:"recommend_tag" xorm:"not null default 0 comment('推荐id 1 低价推荐') INT(10)"`
	AppointmentType int       `json:"appointment_type" xorm:"not null default 0 comment('0 场馆预约 1 私教课程预约 2 大课预约') TINYINT(4)"`
	WeekNum         int       `json:"week_num" xorm:"not null comment('1 周一 2 周二 3 周三 4 周四 5 周五 6 周六 0 周日') TINYINT(4)"`
	Sortorder       int       `json:"sortorder" xorm:"not null default 0 comment('排序权重 倒序') INT(11)"`
	CoachId         int64     `json:"coach_id" xorm:"not null default 0 comment('老师id') BIGINT(20)"`
	CourseId        int64     `json:"course_id" xorm:"not null default 0 comment('课程id') BIGINT(20)"`
	UnitDuration    int       `json:"unit_duration" xorm:"not null default 0 comment('单位时长 [秒] tips:[用于计费使用] 例：每十五分钟/10元 则该字段为900') INT(11)"`
	UnitPrice       int       `json:"unit_price" xorm:"not null default 0 comment('单位价格 [分] tips:[用于计费使用] 例：每十五分钟/10元 则该字段为1000') INT(11)"`
	StartNode       time.Time `json:"start_node" xorm:"not null default '00:00:00' comment('开始时间节点') TIME"`
	EndNode         time.Time `json:"end_node" xorm:"not null default '00:00:00' comment('结束时间节点') TIME"`
	PeriodNum       int       `json:"period_num" xorm:"not null default 0 comment('总课时数') INT(6)"`
	VenueId         int64     `json:"venue_id" xorm:"not null comment('场馆id') BIGINT(20)"`
	CreateAt        int       `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int       `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
