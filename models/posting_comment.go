package models

type PostingComment struct {
	Id                  int64  `json:"id" xorm:"pk autoincr comment('评论id') BIGINT(20)"`
	UserId              string `json:"user_id" xorm:"not null comment('评论人userId') index VARCHAR(60)"`
	PostId              int64  `json:"post_id" xorm:"not null comment('帖子id') index BIGINT(20)"`
	ParentCommentId     int64  `json:"parent_comment_id" xorm:"not null default 0 comment('父评论id') index BIGINT(20)"`
	ParentCommentUserId string `json:"parent_comment_user_id" xorm:"not null default '' comment('父评论的用户id') VARCHAR(60)"`
	ReplyCommentId      int64  `json:"reply_comment_id" xorm:"not null default 0 comment('被回复的评论id') BIGINT(20)"`
	ReplyCommentUserId  string `json:"reply_comment_user_id" xorm:"default '' comment('被回复的评论用户id') VARCHAR(60)"`
	CommentLevel        int    `json:"comment_level" xorm:"not null default 1 comment('评论等级[ 1 一级评论 默认 ，2 二级评论]') TINYINT(4)"`
	Content             string `json:"content" xorm:"not null default '' comment('评论的内容') VARCHAR(1000)"`
	Status              int    `json:"status" xorm:"not null default 1 comment('状态 (1 有效，0 逻辑删除)') index TINYINT(2)"`
	IsTop               int    `json:"is_top" xorm:"not null default 0 comment('置顶状态[ 1 置顶，0 不置顶 默认 ]') TINYINT(2)"`
	CreateAt            int    `json:"create_at" xorm:"not null default 0 comment('创建时间') index INT(11)"`
}
