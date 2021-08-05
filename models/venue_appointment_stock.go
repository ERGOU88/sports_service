package models

type VenueAppointmentStock struct {
	Date         int64 `json:"date" xorm:"not null pk comment('日期 例如 20211001') BIGINT(20)"`
	QuotaNum     int   `json:"quota_num" xorm:"not null default 0 comment('配额：可容纳人数/可预约人数 -1表示没有限制') INT(10)"`
	PurchasedNum int   `json:"purchased_num" xorm:"not null default 0 comment('已购买数量 [冻结库存]') INT(10)"`
	CreateAt     int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt     int   `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
