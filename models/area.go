package models

type Area struct {
	Code   int64  `json:"code" xorm:"not null pk comment('区划代码') BIGINT(12)"`
	Name   string `json:"name" xorm:"not null default '' comment('名称') index VARCHAR(128)"`
	Level  int    `json:"level" xorm:"not null comment('级别1-5,省市县镇村') index TINYINT(1)"`
	Pcode  int64  `json:"pcode" xorm:"comment('父级区划代码') index BIGINT(12)"`
	IsShow int    `json:"is_show" xorm:"not null default 0 comment('0 展示 1 不展示') TINYINT(1)"`
}
