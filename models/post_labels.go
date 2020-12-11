package models

type PostLabels struct {
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	LabelId   string `json:"label_id" xorm:"not null pk default '' comment('标签id') VARCHAR(521)"`
	LabelName string `json:"label_name" xorm:"comment('标签名') VARCHAR(521)"`
	PostId    int64  `json:"post_id" xorm:"not null pk comment('帖子id') BIGINT(20)"`
}
