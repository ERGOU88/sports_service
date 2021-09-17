package models

type RecommendInfoSection struct {
	Id          int    `json:"id" xorm:"not null pk autoincr comment('板块id') INT(11)"`
	SectionType int    `json:"section_type" xorm:"not null default 0 comment('0首页板块') TINYINT(1)"`
	Sortorder   int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Name        string `json:"name" xorm:"comment('板块名称') VARCHAR(60)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('0正常 1废弃') TINYINT(1)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
}
