package models

type SocialAccountLogin struct {
	UserId     string `json:"user_id" xorm:"not null pk comment('用户id') VARCHAR(60)"`
	Unionid    string `json:"unionid" xorm:"not null default '' comment('社交平台关联id') VARCHAR(256)"`
	SocialType int    `json:"social_type" xorm:"not null pk default 0 comment('区分社交软件 1 微信关联id 2 微博关联id 3 qq关联id 4 微信小程序') TINYINT(2)"`
	Status     int    `json:"status" xorm:"default 0 comment('0 正常 1 封禁') TINYINT(1)"`
	CreateAt   int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	OpenId     string `json:"open_id" xorm:"not null default '' comment('社交平台 用户唯一标识') VARCHAR(256)"`
}
