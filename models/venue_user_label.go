package models

type VenueUserLabel struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	Date      int64  `json:"date" xorm:"not null comment('日期 例如 20211001') BIGINT(20)"`
	TimeNode  string `json:"time_node" xorm:"not null default '' comment('时间节点 例如 10:00-12:00') VARCHAR(128)"`
	UserId    string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	LabelType int    `json:"label_type" xorm:"not null default 0 comment('0为用户添加标签 1为系统添加标签') TINYINT(1)"`
	LabelId   int64  `json:"label_id" xorm:"not null comment('标签id') BIGINT(20)"`
	LabelName string `json:"label_name" xorm:"not null comment('标签名') VARCHAR(60)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 有效 1 废弃') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	VenueId   int64  `json:"venue_id" xorm:"not null comment('场馆id') BIGINT(20)"`
}
