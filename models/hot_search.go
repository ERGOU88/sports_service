package models

type HotSearch struct {
	CreateAt         int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	HotSearchContent string `json:"hot_search_content" xorm:"not null comment('热门搜索内容 多个用逗号分隔 例如：FPV,电竞,无人机') VARCHAR(128)"`
	Id               int    `json:"id" xorm:"not null pk autoincr comment('自增主键') INT(11)"`
	UpdateAt         int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
