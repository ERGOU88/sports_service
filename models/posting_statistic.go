package models

type PostingStatistic struct {
	PostingId   int64 `json:"posting_id" xorm:"not null pk comment('帖子id') BIGINT(20)"`
	FabulousNum int   `json:"fabulous_num" xorm:"not null default 0 comment('点赞数') INT(11)"`
	BrowseNum   int   `json:"browse_num" xorm:"not null default 0 comment('浏览数') INT(11)"`
	ShareNum    int   `json:"share_num" xorm:"not null default 0 comment('分享/转发数') INT(11)"`
	RewardNum   int   `json:"reward_num" xorm:"not null default 0 comment('打赏的游币数') INT(11)"`
	CommentNum  int   `json:"comment_num" xorm:"not null default 0 comment('评论数') INT(11)"`
	CollectNum  int   `json:"collect_num" xorm:"not null default 0 comment('收藏数') INT(11)"`
	CreateAt    int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int   `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
