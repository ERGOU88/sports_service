package models

type CollectPostRecord struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	PostId   int64  `json:"post_id" xorm:"not null comment('帖子id') index BIGINT(20)"`
	Status   int    `json:"status" xorm:"not null comment('1 收藏 2 取消收藏') TINYINT(1)"`
	CreateAt int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
