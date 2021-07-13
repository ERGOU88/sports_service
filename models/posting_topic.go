package models

type PostingTopic struct {
	PostingId int64  `json:"posting_id" xorm:"not null pk comment('帖子id') BIGINT(20)"`
	TopicId   int    `json:"topic_id" xorm:"not null pk comment('话题id') INT(11)"`
	TopicName string `json:"topic_name" xorm:"comment('话题名') VARCHAR(521)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('帖子审核通过 则status为1 其他情况默认为0') TINYINT(1)"`
}
