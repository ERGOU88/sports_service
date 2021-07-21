package models

type Links struct {
	Id        int    `json:"id" xorm:"not null pk autoincr comment('自增id') INT(11)"`
	Url       string `json:"url" xorm:"not null comment('长连接') unique VARCHAR(200)"`
	Keyword   string `json:"keyword" xorm:"not null comment('短链接码') unique VARCHAR(100)"`
	Status    int    `json:"status" xorm:"not null default 1 comment('1系统分配 2用户自定义') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	LinksType int    `json:"links_type" xorm:"not null default 0 comment('0 视频') TINYINT(1)"`
}
