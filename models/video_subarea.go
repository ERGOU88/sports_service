package models

type VideoSubarea struct {
	Id          int    `json:"id" xorm:"not null pk autoincr comment('分区id') INT(11)"`
	Sortorder   int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	SubareaName string `json:"subarea_name" xorm:"comment('分区名') VARCHAR(60)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('0正常 1废弃') TINYINT(1)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	SysId       int    `json:"sys_id" xorm:"not null default 0 comment('系统账号id') INT(11)"`
	SysUser     string `json:"sys_user" xorm:"not null default '' comment('系统用户昵称') VARCHAR(128)"`
}
