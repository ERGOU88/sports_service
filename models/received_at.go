package models

type ReceivedAt struct {
	Id           int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	ToUserId     string `json:"to_user_id" xorm:"not null comment('被@的用户id') index VARCHAR(60)"`
	UserId       string `json:"user_id" xorm:"not null comment('执行@的用户id') index VARCHAR(60)"`
	ComposeId    int64  `json:"compose_id" xorm:"not null comment('视频id/帖子id/视频评论id/帖子评论id') BIGINT(20)"`
	TopicType    int    `json:"topic_type" xorm:"not null comment('1.视频评论、回复中@ 2.帖子评论、回复中@ 3.视频评论/回复 4.帖子评论/回复 5.发布帖子时候@的用户') TINYINT(2)"`
	CreateAt     int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	CommentLevel int    `json:"comment_level" xorm:"not null default 1 comment('评论等级[ 1 一级评论 默认 ，2 二级评论]') TINYINT(4)"`
	Status       int    `json:"status" xorm:"not null default 1 comment('@状态 1 正常 0 作品待审核') TINYINT(1)"`
	UpdateAt     int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
