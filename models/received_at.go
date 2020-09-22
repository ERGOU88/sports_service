package models

type ReceivedAt struct {
	CommentId int64  `json:"comment_id" xorm:"not null comment('评论id') BIGINT(20)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Id        int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	ToUserId  string `json:"to_user_id" xorm:"not null comment('被@的用户id') index VARCHAR(60)"`
	TopicType int    `json:"topic_type" xorm:"not null comment('1.视频 2.帖子 3.评论') TINYINT(2)"`
	UserId    string `json:"user_id" xorm:"not null comment('执行@的用户id') index VARCHAR(60)"`
}
