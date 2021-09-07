package models

type Banner struct {
	Id        int    `json:"id" xorm:"not null pk autoincr comment('主键') INT(10)"`
	Title     string `json:"title" xorm:"not null default '' comment('标题') VARCHAR(255)"`
	Cover     string `json:"cover" xorm:"not null default '' comment('banner封面') VARCHAR(512)"`
	Explain   string `json:"explain" xorm:"not null default '' comment('说明') VARCHAR(255)"`
	JumpUrl   string `json:"jump_url" xorm:"not null default '' comment('跳转地址') VARCHAR(512)"`
	ShareUrl  string `json:"share_url" xorm:"not null default '' comment('分享地址') VARCHAR(255)"`
	Type      int    `json:"type" xorm:"not null default 1 comment('1 首页 2 赛事 3 官网banner') INT(1)"`
	StartTime int    `json:"start_time" xorm:"not null default 0 comment('上架时间') INT(11)"`
	EndTime   int    `json:"end_time" xorm:"not null default 0 comment('下架时间') INT(11)"`
	Sortorder int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0待上架 1上架 2 已过期') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	JumpType  int    `json:"jump_type" xorm:"not null comment('跳转类型 0 站内跳转 1 站外跳转') TINYINT(1)"`
}
