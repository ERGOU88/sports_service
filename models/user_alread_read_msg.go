package models

type UserAlreadReadMsg struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	SystemId int64  `json:"system_id" xorm:"not null comment('系统消息ID') index BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null default '' comment('用户id') index VARCHAR(60)"`
	CreateAt int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
