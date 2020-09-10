package models

type UserBrowseVideo struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	VideoId  int64  `json:"video_id" xorm:"not null comment('视频id') BIGINT(20)"`
	CreateAt int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
