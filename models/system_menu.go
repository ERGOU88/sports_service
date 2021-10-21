package models

import (
	"time"
)

type SystemMenu struct {
	MenuId   int       `json:"menu_id" xorm:"not null pk autoincr comment('菜单id') INT(11)"`
	MenuName string    `json:"menu_name" xorm:"not null default '' comment('菜单名称 例如 Role') VARCHAR(128)"`
	Title    string    `json:"title" xorm:"not null default '' comment('抬头 例 角色管理') VARCHAR(128)"`
	Icon     string    `json:"icon" xorm:"not null default '' comment('图标') VARCHAR(128)"`
	Path     string    `json:"path" xorm:"not null comment('访问路径') VARCHAR(128)"`
	Action   string    `json:"action" xorm:"not null default '' comment('请求方式: POST/GET/DELETE等') VARCHAR(16)"`
	ParentId int       `json:"parent_id" xorm:"not null default 0 comment('父级id') INT(11)"`
	Sort     int       `json:"sort" xorm:"not null default 0 comment('排序权重') INT(11)"`
	CreateBy int       `json:"create_by" xorm:"not null default 0 comment('创建者id') INT(11)"`
	UpdateBy int       `json:"update_by" xorm:"not null default 0 comment('更新者id') INT(11)"`
	CreateAt time.Time `json:"create_at" xorm:"TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" xorm:"TIMESTAMP"`
}
