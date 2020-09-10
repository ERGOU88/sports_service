package models

type PostComment struct {
	CommentId int64  `json:"comment_id" xorm:"not null pk autoincr comment('评论id') BIGINT(20)"`
	PostId    int64  `json:"post_id" xorm:"not null comment('帖子id') index BIGINT(20)"`
	Content   string `json:"content" xorm:"not null comment('内容') VARCHAR(1024)"`
	FromUid   string `json:"from_uid" xorm:"not null comment('评论内容的用户id') index VARCHAR(60)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0展示 1不展示') TINYINT(1)"`
}
