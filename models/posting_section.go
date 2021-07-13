package models

type PostingSection struct {
	PostingId   int64  `json:"posting_id" xorm:"not null pk comment('帖子id') BIGINT(20)"`
	SectionId   int    `json:"section_id" xorm:"not null pk comment('版块id') INT(11)"`
	SectionName string `json:"section_name" xorm:"comment('版块名') VARCHAR(521)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	SectionType int    `json:"section_type" xorm:"not null comment('0系统分配板块 1用户自定义板块') TINYINT(1)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('帖子审核通过 则status为1 其他情况默认为0') TINYINT(1)"`
}
