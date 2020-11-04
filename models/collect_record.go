package models

type CollectRecord struct {
	ComposeId   int64  `json:"compose_id" xorm:"not null comment('作品id') index BIGINT(20)"`
	ComposeType int    `json:"compose_type" xorm:"not null default 0 comment('0 视频 1 帖子') TINYINT(2)"`
	CreateAt    int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	Id          int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	Status      int    `json:"status" xorm:"not null comment('1 收藏 0 取消收藏') TINYINT(1)"`
	ToUserId    string `json:"to_user_id" xorm:"not null comment('作品发布者用户id') index VARCHAR(60)"`
	UpdateAt    int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
	UserId      string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
}
