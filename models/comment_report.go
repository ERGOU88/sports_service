package models

type CommentReport struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UserId    string `json:"user_id" xorm:"not null default '' comment('举报人用户id') VARCHAR(60)"`
	CommentId int64  `json:"comment_id" xorm:"not null comment('评论id') index BIGINT(20)"`
	Reason    string `json:"reason" xorm:"not null comment('举报理由') VARCHAR(300)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
}
