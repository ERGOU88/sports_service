package models

type VideoLabels struct {
	VideoId   int64  `json:"video_id" xorm:"not null pk comment('视频id') BIGINT(20)"`
	LabelId   string `json:"label_id" xorm:"not null pk default '' comment('标签id') VARCHAR(521)"`
	LabelName string `json:"label_name" xorm:"comment('标签名') VARCHAR(521)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('视频审核通过 则status为1 其他情况默认为0') TINYINT(1)"`
}
