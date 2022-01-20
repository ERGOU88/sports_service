package models

type ShopServiceConf struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	Service  string `json:"service" xorm:"not null comment('服务') VARCHAR(60)"`
	Icon     string `json:"icon" xorm:"not null default '' comment('图标') VARCHAR(256)"`
	Status   int    `json:"status" xorm:"not null default 0 comment('0可用 1不可用') TINYINT(1)"`
	CreateAt int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	Describe string `json:"describe" xorm:"not null default '' comment('服务描述') VARCHAR(512)"`
}
