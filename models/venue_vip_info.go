package models

type VenueVipInfo struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null comment('用户id') unique VARCHAR(60)"`
	Level    int    `json:"level" xorm:"not null default 0 comment('会员等级 预留字段') INT(2)"`
	VenueId  int64  `json:"venue_id" xorm:"not null comment('场馆id') BIGINT(20)"`
	StartTm  int64  `json:"start_tm" xorm:"not null comment('会员开始时间戳') BIGINT(20)"`
	EndTm    int64  `json:"end_tm" xorm:"not null comment('会员结束时间戳') BIGINT(20)"`
	Duration int64  `json:"duration" xorm:"not null comment('会员在场馆内可用时长') BIGINT(20)"`
	CreateAt int    `json:"create_at" xorm:"not null default 0 INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"not null default 0 INT(11)"`
}
