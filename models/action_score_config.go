package models

type ActionScoreConfig struct {
	Id            int64 `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	WorkType      int   `json:"work_type" xorm:"not null comment('1 视频 2 帖子 3 资讯') TINYINT(2)"`
	FabulousScore int   `json:"fabulous_score" xorm:"not null default 0 comment('点赞得分') INT(11)"`
	BrowseScore   int   `json:"browse_score" xorm:"not null default 0 comment('浏览得分') INT(11)"`
	ShareScore    int   `json:"share_score" xorm:"not null default 0 comment('分享得分') INT(11)"`
	RewardScore   int   `json:"reward_score" xorm:"not null default 0 comment('打赏得分') INT(11)"`
	BarrageScore  int   `json:"barrage_score" xorm:"not null default 0 comment('弹幕得分') INT(11)"`
	CommentScore  int   `json:"comment_score" xorm:"not null default 0 comment('评论得分') INT(11)"`
	CollectScore  int   `json:"collect_score" xorm:"not null default 0 comment('收藏得分') INT(11)"`
	TopScore      int   `json:"top_score" xorm:"not null default 0 comment('置顶需达到的分值') INT(11)"`
	CreateAt      int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int   `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
