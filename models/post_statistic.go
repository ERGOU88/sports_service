package models

type PostStatistic struct {
	PostId   int64 `json:"post_id" xorm:"not null pk comment('帖子id') BIGINT(20)"`
	Fabulous int   `json:"fabulous" xorm:"not null default 0 comment('点赞数') INT(11)"`
	Browse   int   `json:"browse" xorm:"not null default 0 comment('浏览数') INT(11)"`
	Share    int   `json:"share" xorm:"not null default 0 comment('分享数') INT(11)"`
	Reward   int   `json:"reward" xorm:"not null default 0 comment('打赏的游币数') INT(11)"`
	CreateAt int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt int   `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
