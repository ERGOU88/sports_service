package models

type SectionConfig struct {
	Id          int    `json:"id" xorm:"not null pk autoincr comment('主键') INT(11)"`
	CircleId    int    `json:"circle_id" xorm:"not null comment('圈子id') index INT(11)"`
	SectionName string `json:"section_name" xorm:"not null comment('板块名称') VARCHAR(100)"`
	Sortorder   int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('0 未操作 1 展示  2 隐藏') TINYINT(1)"`
	CreateAt    int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
