package models

type VideoLabels struct {
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	LabelId   string `json:"label_id" xorm:"not null pk default '' comment('标签id') VARCHAR(521)"`
	LabelName string `json:"label_name" xorm:"comment('标签名') VARCHAR(521)"`
	VideoId   int64  `json:"video_id" xorm:"not null pk comment('视频id') BIGINT(20)"`
}
