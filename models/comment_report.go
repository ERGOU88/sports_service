package models

type CommentReport struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UserId      string `json:"user_id" xorm:"not null default '' comment('举报人用户id') VARCHAR(60)"`
	CommentId   int64  `json:"comment_id" xorm:"not null comment('评论id') index BIGINT(20)"`
	Reason      string `json:"reason" xorm:"not null comment('举报理由') VARCHAR(300)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	CommentType int    `json:"comment_type" xorm:"not null default 1 comment('1视频评论 2帖子评论') TINYINT(1)"`
}
