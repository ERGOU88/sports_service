package models

type UserBrowseRecord struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	UserId      string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	ComposeId   int64  `json:"compose_id" xorm:"not null comment('作品id') BIGINT(20)"`
	ComposeType int    `json:"compose_type" xorm:"not null default 0 comment('0 视频 1 帖子 2 资讯') TINYINT(2)"`
	CreateAt    int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
