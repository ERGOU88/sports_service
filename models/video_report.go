package models

type VideoReport struct {
	CreateAt int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Id       int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null default '' comment('用户id') VARCHAR(60)"`
	VideoId  int64  `json:"video_id" xorm:"not null comment('视频id') index BIGINT(20)"`
}
