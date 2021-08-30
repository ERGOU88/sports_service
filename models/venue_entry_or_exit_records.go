package models

type VenueEntryOrExitRecords struct {
	Id         int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	UserId     string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	ActionType int    `json:"action_type" xorm:"not null comment('动作类型 1 进场 2 出场') TINYINT(1)"`
	VenueName  string `json:"venue_name" xorm:"not null comment('场馆名称') VARCHAR(60)"`
	VenueId    int64  `json:"venue_id" xorm:"not null comment('场馆ID') BIGINT(20)"`
	Status     int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	CreateAt   int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt   int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
