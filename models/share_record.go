package models

type ShareRecord struct {
	ComposeId     int64  `json:"compose_id" xorm:"not null comment('作品id（视频/帖子id）') BIGINT(20)"`
	Content       string `json:"content" xorm:"not null comment('分享的整体内容（json）') TEXT"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Id            int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	SharePlatform int    `json:"share_platform" xorm:"not null comment('分享平台 1 微信 2 微博 3 qq') TINYINT(2)"`
	ShareType     int    `json:"share_type" xorm:"not null comment('分享类型 1 分享视频 2 分享帖子') TINYINT(2)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('0展示 1不展示') TINYINT(1)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
}
