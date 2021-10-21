package models

import (
	"time"
)

type SystemRoleMenu struct {
	RoleId   int       `json:"role_id" xorm:"not null pk comment('角色id') INT(11)"`
	MenuId   int       `json:"menu_id" xorm:"not null pk comment('菜单id') INT(11)"`
	RoleName string    `json:"role_name" xorm:"not null default '' comment('角色名称') VARCHAR(128)"`
	CreateAt time.Time `json:"create_at" xorm:"TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" xorm:"TIMESTAMP"`
}
