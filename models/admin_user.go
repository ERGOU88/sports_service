package models

type AdminUser struct {
	CreateAt   int    `json:"create_at" xorm:"not null comment('注册时间') INT(11)"`
	Id         int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Password   string `json:"password" xorm:"not null comment('登录密码') CHAR(32)"`
	Phone      string `json:"phone" xorm:"default '' comment('手机号') CHAR(11)"`
	Salt       string `json:"salt" xorm:"not null comment('salt') CHAR(10)"`
	SubAccount int64  `json:"sub_account" xorm:"not null default 0 comment('子账号标识：0则为主账号，大于0则为子账号，值等于主账号ID') BIGINT(11)"`
	UpdateAt   int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
	Username   string `json:"username" xorm:"not null comment('用户名') VARCHAR(255)"`
}
