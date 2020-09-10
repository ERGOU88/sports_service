package models

type SearchHistory struct {
	Id            int    `json:"id" xorm:"not null pk autoincr comment('自增主键') INT(11)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	SearchContent string `json:"search_content" xorm:"not null comment('搜索的内容') VARCHAR(128)"`
	CreateAt      int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	Status        int    `json:"status" xorm:"not null default 1 comment('1 正常 2 已删除') TINYINT(1)"`
}
