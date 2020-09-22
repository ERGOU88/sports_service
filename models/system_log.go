package models

type SystemLog struct {
	CreateAt int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Id       int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	LogCont  string `json:"log_cont" xorm:"comment('操作内容') MEDIUMTEXT"`
	LogText  string `json:"log_text" xorm:"comment('备注') MEDIUMTEXT"`
	LogType  string `json:"log_type" xorm:"default '' comment('记录类型') VARCHAR(200)"`
	SysId    int64  `json:"sys_id" xorm:"not null default 0 comment('系统账号ID') BIGINT(20)"`
	SysRole  string `json:"sys_role" xorm:"default '' comment('用户角色') VARCHAR(255)"`
	SysUser  string `json:"sys_user" xorm:"default '' comment('用户昵称') VARCHAR(255)"`
}
