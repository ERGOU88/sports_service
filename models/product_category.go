package models

type ProductCategory struct {
	CategoryId   int    `json:"category_id" xorm:"not null pk autoincr comment('分类id') INT(11)"`
	CategoryName string `json:"category_name" xorm:"not null default '' comment('分类名称') VARCHAR(50)"`
	ShortName    string `json:"short_name" xorm:"not null default '' comment('简称') VARCHAR(50)"`
	Pid          int    `json:"pid" xorm:"not null default 0 comment('分类上级') index(pid_level) INT(11)"`
	Level        int    `json:"level" xorm:"not null default 0 comment('层级') index(pid_level) INT(11)"`
	IsShow       int    `json:"is_show" xorm:"not null default 0 comment('是否显示（0显示  1不显示）') INT(11)"`
	Sortorder    int    `json:"sortorder" xorm:"not null default 0 comment('排序') INT(11)"`
	Image        string `json:"image" xorm:"not null default '' comment('分类图片') VARCHAR(255)"`
	Keywords     string `json:"keywords" xorm:"not null default '' comment('分类页面关键字') VARCHAR(255)"`
	Description  string `json:"description" xorm:"not null default '' comment('分类介绍') VARCHAR(255)"`
	CreateAt     int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt     int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
