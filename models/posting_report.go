package models

type PostingReport struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null default '' comment('用户id') VARCHAR(60)"`
	PostId   int64  `json:"post_id" xorm:"not null comment('帖子id') index BIGINT(20)"`
	Reason   string `json:"reason" xorm:"not null default '' comment('举报理由') VARCHAR(100)"`
	CreateAt int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
}
