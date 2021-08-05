package models

type VenuePersonalLabelConf struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	LabelName string `json:"label_name" xorm:"not null comment('标签名') VARCHAR(60)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0可用 1不可用') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
