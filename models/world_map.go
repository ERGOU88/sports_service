package models

type WorldMap struct {
	Id        int    `json:"id" xorm:"not null pk autoincr comment('主键id') INT(11)"`
	Name      string `json:"name" xorm:"not null comment('国家(省份/城市)名称') VARCHAR(50)"`
	Code      string `json:"code" xorm:"not null comment('国家(省份/城市)编码') CHAR(4)"`
	Pid       int    `json:"pid" xorm:"not null default 0 comment('父级id') index INT(11)"`
	Layer     int    `json:"layer" xorm:"not null default 0 comment('层级') TINYINT(2)"`
	Sortorder int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 展示 1 隐藏') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
