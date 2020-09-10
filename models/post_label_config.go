package models

type PostLabelConfig struct {
	LabelId   int    `json:"label_id" xorm:"not null pk autoincr comment('标签id') INT(11)"`
	Pid       int    `json:"pid" xorm:"not null comment('父类id 0为1级分类') INT(11)"`
	Sortorder int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 1 comment('类别状态1-正常,2-已废弃') INT(1)"`
	LabelName string `json:"label_name" xorm:"not null comment('标签名称') VARCHAR(64)"`
	Icon      string `json:"icon" xorm:"not null default '' comment('标签icon') VARCHAR(256)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
