package models

import (
	"time"
)

type SystemRole struct {
	RoleId    int       `json:"role_id" xorm:"not null pk autoincr comment('角色id') INT(11)"`
	RoleName  string    `json:"role_name" xorm:"not null default '' comment('角色名称') VARCHAR(128)"`
	Status    int       `json:"status" xorm:"not null default 0 comment('0正常 1禁用') TINYINT(1)"`
	RoleKey   string    `json:"role_key" xorm:"not null default '' comment('角色标示') VARCHAR(128)"`
	RoleSort  int       `json:"role_sort" xorm:"not null default 0 comment('权限排序') INT(11)"`
	CreateBy  int       `json:"create_by" xorm:"not null default 0 comment('创建者id') INT(11)"`
	UpdateBy  int       `json:"update_by" xorm:"not null default 0 comment('更新者id') INT(11)"`
	Remark    string    `json:"remark" xorm:"not null default '' comment('备注') VARCHAR(255)"`
	CreatedAt time.Time `json:"created_at" xorm:"TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" xorm:"TIMESTAMP"`
}
