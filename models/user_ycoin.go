package models

type UserYcoin struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	Ycoin    int    `json:"ycoin" xorm:"not null comment('游币数') index INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
