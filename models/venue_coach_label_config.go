package models

type VenueCoachLabelConfig struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	LabelName string `json:"label_name" xorm:"not null comment('标签名称') VARCHAR(30)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 有效 1 废弃') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
