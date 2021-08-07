package models

type VenueAppointmentStock struct {
	Date            int64  `json:"date" xorm:"not null pk comment('日期 例如 20211001') BIGINT(20)"`
	TimeNode        string `json:"time_node" xorm:"not null pk default '' comment('开始时间节点 例如 10:00-12:00') VARCHAR(128)"`
	QuotaNum        int    `json:"quota_num" xorm:"not null default 0 comment('配额：可容纳人数/可预约人数 -1表示没有限制') INT(10)"`
	PurchasedNum    int    `json:"purchased_num" xorm:"not null default 0 comment('已购买数量 [冻结库存]') INT(10)"`
	AppointmentType int    `json:"appointment_type" xorm:"not null default 0 comment('0 场馆预约 1 私教预约 2 课程预约') TINYINT(1)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
