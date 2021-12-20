package models

type Area struct {
	Id        int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Pid       int    `json:"pid" xorm:"not null default 0 comment('父级') index INT(11)"`
	Name      string `json:"name" xorm:"not null default '' comment('名称') index(area) VARCHAR(50)"`
	Shortname string `json:"shortname" xorm:"not null default '' comment('简称') index(area) VARCHAR(30)"`
	Longitude string `json:"longitude" xorm:"not null default '' comment('经度') index(longitude) VARCHAR(30)"`
	Latitude  string `json:"latitude" xorm:"not null default '' comment('纬度') index(longitude) VARCHAR(30)"`
	Level     int    `json:"level" xorm:"not null default 0 comment('级别') index(level) TINYINT(2)"`
	Sortorder int    `json:"sortorder" xorm:"not null default 0 comment('排序') index(level) INT(10)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('状态 0 有效 1 废弃') index(level) TINYINT(1)"`
}
