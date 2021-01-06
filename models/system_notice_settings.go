package models

type SystemNoticeSettings struct {
	AttentionPushSet int    `json:"attention_push_set" xorm:"not null default 0 comment('状态（0：接收推送；1：拒绝推送 关注推送）') TINYINT(1)"`
	CommentPushSet   int    `json:"comment_push_set" xorm:"not null default 0 comment('状态（0：接收推送；1：拒绝推送 包含评论/回复推送）') TINYINT(1)"`
	CreateAt         int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	SharePushSet     int    `json:"share_push_set" xorm:"not null default 0 comment('状态（0：接收推送；1：拒绝推送 包含评论/回复）') TINYINT(1)"`
	SlotPushSet      int    `json:"slot_push_set" xorm:"not null default 0 comment('状态（0：接收推送；1：拒绝推送 投币推送）') TINYINT(1)"`
	ThumbUpPushSet   int    `json:"thumb_up_push_set" xorm:"not null default 0 comment('状态（0：接收推送；1：拒绝推送 点赞推送）') TINYINT(1)"`
	UpdateAt         int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	UserId           string `json:"user_id" xorm:"not null pk comment('用户id') VARCHAR(60)"`
}
