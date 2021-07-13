package models

type ForwardRecord struct {
	Id          int64  `json:"id" xorm:"pk comment('帖子id') BIGINT(20)"`
	UserId      string `json:"user_id" xorm:"not null comment('转发者用户id') index VARCHAR(60)"`
	ToUserId    string `json:"to_user_id" xorm:"not null comment('被转发者用户id') VARCHAR(60)"`
	ForwardType int    `json:"forward_type" xorm:"not null default 0 comment('转发类型 0 视频 1 贴子') TINYINT(1)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('状态（0正常 1删除）') TINYINT(1)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
