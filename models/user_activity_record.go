package models

type UserActivityRecord struct {
	Id           int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	UserId       string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	CreateAt     int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt     int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	ActivityType int    `json:"activity_type" xorm:"not null default 0 comment('用户活跃类型') INT(6)"`
}
