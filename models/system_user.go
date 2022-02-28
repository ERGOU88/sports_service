package models

import (
	"time"
)

type SystemUser struct {
	UserId      int       `json:"user_id" xorm:"not null pk autoincr comment('管理员id') INT(11)"`
	NickName    string    `json:"nick_name" xorm:"not null default '' comment('昵称/姓名') VARCHAR(128)"`
	Phone       string    `json:"phone" xorm:"not null default '' comment('手机号') VARCHAR(11)"`
	RoleId      int       `json:"role_id" xorm:"not null default 0 comment('角色id') INT(11)"`
	Salt        string    `json:"salt" xorm:"not null default '' comment('盐') VARCHAR(255)"`
	Avatar      string    `json:"avatar" xorm:"not null default '' comment('头像') VARCHAR(255)"`
	Sex         string    `json:"sex" xorm:"not null default '' comment('性别') VARCHAR(255)"`
	Email       string    `json:"email" xorm:"not null default '' comment('邮箱') VARCHAR(128)"`
	DeptId      int       `json:"dept_id" xorm:"not null default 0 comment('部门id 预留') INT(11)"`
	Username    string    `json:"username" xorm:"not null comment('账号') VARCHAR(64)"`
	Password    string    `json:"password" xorm:"not null comment('密码') VARCHAR(128)"`
	CreateBy    int       `json:"create_by" xorm:"not null default 0 comment('创建者id') INT(11)"`
	UpdateBy    int       `json:"update_by" xorm:"not null default 0 comment('更新者id') INT(11)"`
	Remark      string    `json:"remark" xorm:"not null default '' comment('备注') VARCHAR(255)"`
	Status      int       `json:"status" xorm:"not null default 1 comment('1 正常 -1 封禁') TINYINT(1)"`
	CreateAt    time.Time `json:"create_at" xorm:"TIMESTAMP"`
	UpdateAt    time.Time `json:"update_at" xorm:"TIMESTAMP"`
	AccountType int       `json:"account_type" xorm:"not null default 0 comment('0 域账号 1 自定义账号') TINYINT(1)"`
}
