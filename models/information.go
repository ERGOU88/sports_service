package models

type Information struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('资讯id') BIGINT(20)"`
	RelatedId int64  `json:"related_id" xorm:"not null default 0 comment('关联ID (包含直播id /  视频首页板块id)') BIGINT(20)"`
	Cover     string `json:"cover" xorm:"not null default '' comment('封面') VARCHAR(256)"`
	Title     string `json:"title" xorm:"comment('标题') MEDIUMTEXT"`
	Content   string `json:"content" xorm:"comment('内容') MEDIUMTEXT"`
	JumpUrl   string `json:"jump_url" xorm:"not null default '' comment('跳转地址') VARCHAR(500)"`
	Sortorder int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 正常 1 隐藏') TINYINT(1)"`
	UserId    string `json:"user_id" xorm:"not null default '' comment('官方账号 用户id') VARCHAR(60)"`
	PubType   int    `json:"pub_type" xorm:"not null default 1 comment('后台发布类型 1. 发布至赛事模块 2. 发布至视频首页板块') TINYINT(2)"`
}
