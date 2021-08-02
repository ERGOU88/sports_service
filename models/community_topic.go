package models

type CommunityTopic struct {
	Id        int    `json:"id" xorm:"not null pk autoincr comment('主键') INT(11)"`
	TopicName string `json:"topic_name" xorm:"not null comment('话题名称') VARCHAR(100)"`
	Sortorder int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 未操作 1 展示  2 隐藏') TINYINT(1)"`
	IsHot     int    `json:"is_hot" xorm:"not null default 0 comment('是否热门话题 1 热门') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
	Cover     string `json:"cover" xorm:"not null default '' comment('话题封面') VARCHAR(256)"`
	Describe  string `json:"describe" xorm:"not null default '' comment('话题描述') VARCHAR(1000)"`
	SectionId int    `json:"section_id" xorm:"not null default 0 comment('所属板块id') INT(11)"`
}
