package models

type Information struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('资讯id') BIGINT(20)"`
	LiveId    int64  `json:"live_id" xorm:"not null default 0 comment('关联直播（video_live)ID ') BIGINT(20)"`
	Cover     string `json:"cover" xorm:"not null default '' comment('封面') VARCHAR(256)"`
	Title     string `json:"title" xorm:"comment('标题') MEDIUMTEXT"`
	Content   string `json:"content" xorm:"comment('内容') MEDIUMTEXT"`
	JumpUrl   string `json:"jump_url" xorm:"not null default '' comment('跳转地址') VARCHAR(500)"`
	Sortorder int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UserId    string `json:"user_id" xorm:"not null comment('官方账号 用户id') VARCHAR(60)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 隐藏 1 展示') TINYINT(1)"`
}
