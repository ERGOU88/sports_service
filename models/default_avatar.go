package models

type DefaultAvatar struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	Avatar    string `json:"avatar" xorm:"not null comment('头像地址') VARCHAR(128)"`
	Sortorder int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0展示 1不展示') TINYINT(1)"`
}
